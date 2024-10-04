import {
  BadRequestException,
  Controller,
  ForbiddenException,
  Headers,
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
import { BillingUrlEntity } from "./entities/billing-url.entity";
import { StripeService } from "./stripe.service";

@Roles("ADMIN")
@ApiBearerAuth()
@ApiTags("Organization - Billing")
@Controller()
export class StripeController {
  private readonly logger: Logger = new Logger("StripeController");

  constructor(
    private readonly stripeService: StripeService,
    private organizationsService: OrganizationsService,
    private readonly configService: ConfigService
  ) {}

  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Successfully created url",
    status: 201,
    type: BillingUrlEntity,
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @ApiOperation({
    description:
      "This endpoint will create a billing for an organization to edit their subscription and billing information. Only available on archesai.com. ADMIN ONLY.",
    summary: "Create a billing portal for an organization",
  })
  @Post("/organizations/:orgname/billing")
  async createBillingPortal(
    @Param("orgname") orgName: string
  ): Promise<BillingUrlEntity> {
    if (this.configService.get("FEATURE_BILLING") == false) {
      throw new ForbiddenException("Billing is disabled");
    }
    const organization = await this.organizationsService.findOneByName(orgName);

    return new BillingUrlEntity(
      await this.stripeService.createBillingPortal(
        organization.stripeCustomerId
      )
    );
  }

  @ApiResponse({ description: "Unauthorized", status: 401 })
  @ApiResponse({ description: "Bad Product", status: 400 })
  @ApiResponse({ description: "Not Found", status: 404 })
  @ApiResponse({
    description: "Successfully created url",
    status: 201,
    type: BillingUrlEntity,
  })
  @ApiOperation({
    description:
      "This endpoint will create a checkout session for an organization to purchase a subscription. Only available on archesai.com. ADMIN ONLY.",
    summary: "Create a checkout session for an organization",
  })
  @ApiResponse({ description: "Forbidden", status: 403 })
  @Post("/organizations/:orgname/checkout")
  async createCheckoutSession(
    @Param("orgname") orgName: string,
    @Query("product") product: string
  ) {
    if (this.configService.get("FEATURE_BILLING") == false) {
      throw new ForbiddenException("Billing is disabled");
    }
    const organization = await this.organizationsService.findOneByName(orgName);

    let priceId = "";
    switch (product) {
      case "API":
        priceId = this.configService.get("STRIPE_API_PRICE_ID");
        break;
      case "API_CREDITS":
        priceId = this.configService.get("STRIPE_API_CREDITS_PRICE_ID");
        break;
      case "BASIC":
        priceId = this.configService.get("STRIPE_BASIC_PRICE_ID");
        break;
      default:
        throw new BadRequestException("Invalid product");
    }
    return new BillingUrlEntity(
      await this.stripeService.createCheckoutSession(
        organization.stripeCustomerId,
        {
          price: priceId,
          quantity: 1,
        },
        product == "API_CREDITS"
      )
    );
  }

  @IsPublic()
  @ApiExcludeEndpoint()
  @Post("/webhooks/stripe")
  async handleIncomingEvents(
    @Headers("stripe-signature") signature: string,
    @Req() req: RawBodyRequest<Request>
  ) {
    if (!signature) {
      throw new BadRequestException("Missing stripe-signature header");
    }

    const event = await this.stripeService.constructEventFromPayload(
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
          const price = await this.stripeService.getPrice(lineItem.price.id);
          const credits = price.metadata["credits"];
          await this.organizationsService.addCredits(
            organization.orgname,
            Number(credits)
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

      const subscriptionMap: {
        [key: string]: PlanType;
      } = {
        [this.configService.get("STRIPE_API_PRICE_ID") as string]: "API",
        [this.configService.get("STRIPE_BASIC_PRICE_ID") as string]: "BASIC",
        [this.configService.get("STRIPE_PREMIUM_PRICE_ID") as string]:
          "PREMIUM",
      };
      if (data.status == "active") {
        await this.organizationsService.setPlan(
          organization.orgname,
          subscriptionMap[data.items.data[0].price.id]
        );
      } else {
        await this.organizationsService.setPlan(organization.orgname, "FREE");
      }
    }
  }
}
