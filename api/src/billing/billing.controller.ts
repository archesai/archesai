import {
  BadRequestException,
  Controller,
  Delete,
  ForbiddenException,
  Get,
  Headers,
  NotFoundException,
  Param,
  Post,
  Query,
  RawBodyRequest,
  Req,
} from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import {
  ApiBearerAuth,
  ApiExcludeEndpoint,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from "@nestjs/swagger";
import { PlanType } from "@prisma/client";
import { Stripe } from "stripe";

import { IsPublic } from "../auth/decorators/is-public.decorator";
import { Roles } from "../auth/decorators/roles.decorator";
import { OrganizationsService } from "../organizations/organizations.service";
import { BillingService } from "./billing.service";
import { BillingUrlEntity } from "./entities/billing-url.entity";
import { PaymentMethodEntity } from "./entities/payment-method.entity";
import { PlanEntity } from "./entities/plan.entity";

@Roles("ADMIN")
@ApiBearerAuth()
@ApiTags("Organization - Billing")
@Controller()
export class BillingController {
  private readonly logger: Logger = new Logger("BillingController");

  constructor(
    private readonly billingService: BillingService,
    private organizationsService: OrganizationsService,
    private readonly configService: ConfigService
  ) {}

  @ApiOperation({
    description:
      "This endpoint will create a billing portal for an organization to edit their subscription and billing information. Only available on archesai.com. ADMIN ONLY.",
    summary: "Create a billing portal for an organization",
  })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Successfully created URL",
    status: 201,
    type: BillingUrlEntity,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Post("/organizations/:orgname/billing/portal")
  async createBillingPortal(
    @Param("orgname") orgname: string
  ): Promise<BillingUrlEntity> {
    if (this.configService.get("FEATURE_BILLING") == false) {
      throw new ForbiddenException("Billing is disabled");
    }
    const organization = await this.organizationsService.findOneByName(orgname);

    return new BillingUrlEntity(
      await this.billingService.createBillingPortal(
        organization.stripeCustomerId
      )
    );
  }

  @ApiOperation({
    description:
      "This endpoint will create a checkout session for an organization to purchase a subscription or one-time product. Only available on archesai.com. ADMIN ONLY.",
    summary: "Create a checkout session for an organization",
  })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Bad Request", status: 400 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Successfully created checkout session URL",
    status: 201,
    type: BillingUrlEntity,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Post("/organizations/:orgname/billing/checkout")
  async createCheckoutSession(
    @Param("orgname") orgname: string,
    @Query("planId") planId: string
  ): Promise<BillingUrlEntity> {
    if (this.configService.get("FEATURE_BILLING") == false) {
      throw new ForbiddenException("Billing is disabled");
    }
    const organization = await this.organizationsService.findOneByName(orgname);

    const plans = await this.billingService.listPlans();
    const plan = plans.find((p) => p.id === planId);

    if (!plan) {
      throw new BadRequestException("Invalid plan");
    }

    const priceId = plan.priceId;

    return new BillingUrlEntity(
      await this.billingService.createCheckoutSession(
        organization.stripeCustomerId,
        {
          price: priceId,
          quantity: 1,
        },
        !plan?.recurring?.interval
      )
    );
  }

  @ApiTags("Plans")
  @ApiOperation({
    description: "Get a list of available billing plans",
    summary: "List billing plans",
  })
  @ApiResponse({
    description: "List of plans",
    status: 200,
    type: [PlanEntity],
  })
  @IsPublic()
  @Get("/plans")
  async getPlans(): Promise<PlanEntity[]> {
    const plans = await this.billingService.listPlans();
    return plans.map((plan) => new PlanEntity(plan));
  }

  @ApiOperation({
    description: "List payment methods for an organization",
    summary: "List payment methods",
  })
  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "List of payment methods",
    status: 200,
    type: [PaymentMethodEntity],
  })
  @Get("/organizations/:orgname/billing/payment-methods")
  async listPaymentMethods(@Param("orgname") orgname: string) {
    const organization = await this.organizationsService.findOneByName(orgname);
    const paymentMethods = await this.billingService.listPaymentMethods(
      organization.stripeCustomerId
    );
    return paymentMethods.data;
  }

  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiOperation({
    description: "Remove a payment method from an organization",
    summary: "Remove payment method",
  })
  @Delete("/organizations/:orgname/billing/payment-methods/:paymentMethodId")
  async removePaymentMethod(
    @Param("orgname") orgname: string,
    @Param("paymentMethodId") paymentMethodId: string
  ) {
    const organization = await this.organizationsService.findOneByName(orgname);
    const paymentMethods = await this.billingService.listPaymentMethods(
      organization.stripeCustomerId
    );
    const paymentMethod = paymentMethods.data.find(
      (pm) => pm.id === paymentMethodId
    );
    if (!paymentMethod) {
      throw new NotFoundException("Payment method not found");
    }
    await this.billingService.detachPaymentMethod(paymentMethodId);
    return { success: true };
  }

  @ApiExcludeEndpoint()
  @IsPublic()
  @Post("/webhooks/stripe")
  async stripe_handleIncomingEvents(
    @Headers("stripe-signature") signature: string,
    @Req() req: RawBodyRequest<Request>
  ) {
    if (!signature) {
      throw new BadRequestException("Missing stripe-signature header");
    }

    const event = await this.billingService.constructEventFromPayload(
      signature,
      req.rawBody
    );

    if (event.type == "invoice.paid") {
      const data = event.data.object as Stripe.Invoice;
      if (data.amount_paid > 0) {
        const customerId = data.customer as string;
        const organization =
          await this.organizationsService.findOneByCustomerId(customerId);
        for (const lineItem of data.lines.data) {
          const price = await this.billingService.getPrice(lineItem.price.id);
          const credits = price.metadata["credits"];
          const quantity = lineItem.quantity || 1;
          await this.organizationsService.addCredits(
            organization.orgname,
            Number(credits) * quantity
          );
        }
      }
    }

    if (
      event.type == "customer.subscription.created" ||
      event.type == "customer.subscription.updated" ||
      event.type == "customer.subscription.deleted"
    ) {
      const data = event.data.object as Stripe.Subscription;
      const customerId = data.customer as string;
      const organization =
        await this.organizationsService.findOneByCustomerId(customerId);

      const priceId = data.items.data[0].price.id;
      const price = await this.billingService.getPrice(priceId);
      const planType = price.metadata["plan_type"] as PlanType;

      if (data.status == "active") {
        await this.organizationsService.setPlan(organization.orgname, planType);
      } else {
        await this.organizationsService.setPlan(organization.orgname, "FREE");
      }
    }
  }
}
