import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import Stripe from "stripe";

@Injectable()
export class StripeService {
  private stripe: Stripe;
  constructor(private configService: ConfigService) {
    this.stripe = new Stripe(this.configService.get("STRIPE_PRIVATE_API_KEY"), {
      apiVersion: "2024-06-20",
    });
  }

  async constructEventFromPayload(signature: string, payload: Buffer) {
    const webhookSecret = this.configService.get("STRIPE_WEBHOOK_SECRET");

    return this.stripe.webhooks.constructEvent(
      payload,
      signature,
      webhookSecret,
    );
  }

  async createBillingPortal(customerId: string) {
    return this.stripe.billingPortal.sessions.create({
      customer: customerId,
      return_url: `${this.configService.get("FRONTEND_HOST")}/documents`,
    });
  }

  async createCheckoutSession(
    customerId: string,
    lineItem: { price: string; quantity: number },
    adjustable_quantity: boolean,
  ) {
    const session = await this.stripe.checkout.sessions.create({
      payment_method_types: ["card"],
      ...(adjustable_quantity
        ? {
            invoice_creation: {
              enabled: true,
            },
          }
        : {}),
      allow_promotion_codes: adjustable_quantity ? undefined : true,
      cancel_url: `${this.configService.get("FRONTEND_HOST")}/documents`,
      customer: customerId,
      line_items: [
        {
          ...lineItem,
          adjustable_quantity: adjustable_quantity
            ? { enabled: true, minimum: 1 }
            : undefined,
        },
      ],
      mode: adjustable_quantity ? "payment" : "subscription",
      success_url: `${this.configService.get("FRONTEND_HOST")}/documents`,
    });

    return { url: session.url };
  }

  async createCustomer(name: string, billingEmail: string) {
    return this.stripe.customers.create({
      email: billingEmail,
      name,
    });
  }

  async getCustomer(stripeCustomerId: string) {
    return this.stripe.customers.retrieve(stripeCustomerId);
  }

  async getPrice(id: string) {
    return this.stripe.prices.retrieve(id);
  }

  public async listAllSubscriptions() {
    return this.stripe.subscriptions.list();
  }
}
