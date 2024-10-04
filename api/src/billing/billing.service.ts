// billing.service.ts

import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import Stripe from "stripe";

@Injectable()
export class BillingService {
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
      webhookSecret
    );
  }

  async createBillingPortal(customerId: string) {
    return this.stripe.billingPortal.sessions.create({
      customer: customerId,
      return_url: `${this.configService.get("FRONTEND_HOST")}/settings/organization/billing`,
    });
  }

  async createCheckoutSession(
    customerId: string,
    lineItem: { price: string; quantity: number },
    isOneTime: boolean
  ) {
    const session = await this.stripe.checkout.sessions.create({
      allow_promotion_codes: isOneTime ? undefined : true,
      cancel_url: `${this.configService.get("FRONTEND_HOST")}/settings/organization/billing`,
      customer: customerId,
      invoice_creation: isOneTime
        ? {
            enabled: true,
          }
        : undefined,
      line_items: [
        {
          ...lineItem,
          adjustable_quantity: isOneTime
            ? { enabled: true, minimum: 1 }
            : undefined,
        },
      ],
      mode: isOneTime ? "payment" : "subscription",
      payment_method_types: ["card"],
      success_url: `${this.configService.get("FRONTEND_HOST")}/settings/organization/billing`,
    });

    return { url: session.url };
  }

  async createCustomer(name: string, billingEmail: string) {
    return this.stripe.customers.create({
      email: billingEmail,
      name,
    });
  }

  async createSetupIntent(customerId: string) {
    return this.stripe.setupIntents.create({
      customer: customerId,
      payment_method_types: ["card"],
    });
  }

  async detachPaymentMethod(paymentMethodId: string) {
    return this.stripe.paymentMethods.detach(paymentMethodId);
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

  async listPaymentMethods(customerId: string) {
    return this.stripe.paymentMethods.list({
      customer: customerId,
      type: "card",
    });
  }

  async listPlans() {
    const products = await this.stripe.products.list({
      active: true,
      expand: ["data.default_price"],
    });

    const plans = products.data
      .map((product) => {
        const price = product.default_price as Stripe.Price;
        if (!price) {
          return null;
        }
        return {
          currency: price.currency,
          description: product.description,
          id: product.id,
          metadata: product.metadata,
          name: product.name,
          priceId: price.id,
          priceMetadata: price.metadata,
          recurring: price.recurring,
          unitAmount: price.unit_amount,
        };
      })
      .filter((val) => val !== null);

    console.log(plans);
    return plans;
  }
}
