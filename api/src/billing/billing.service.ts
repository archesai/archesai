import { BadRequestException, Injectable } from "@nestjs/common";
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

  async cancelSubscription(customerId: string) {
    // Retrieve the customer's active subscriptions
    const subscriptions = await this.stripe.subscriptions.list({
      customer: customerId,
      status: "active",
    });

    if (subscriptions.data.length === 0) {
      throw new BadRequestException(
        "No active subscriptions found for this customer."
      );
    }

    // Assuming there is only one subscription per customer
    const subscription = subscriptions.data[0];

    // Cancel the subscription immediately
    await this.stripe.subscriptions.cancel(subscription.id);

    // Alternatively, to cancel at period end, use:
    // await this.stripe.subscriptions.update(subscription.id, {
    //   cancel_at_period_end: true,
    // });
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
      return_url: `${this.configService.get(
        "FRONTEND_HOST"
      )}/settings/organization/billing`,
    });
  }

  async createCheckoutSession(
    customerId: string,
    lineItem: { price: string; quantity: number },
    isOneTime: boolean
  ) {
    const session = await this.stripe.checkout.sessions.create({
      allow_promotion_codes: isOneTime ? undefined : true,
      cancel_url: `${this.configService.get(
        "FRONTEND_HOST"
      )}/settings/organization/billing`,
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
      success_url: `${this.configService.get(
        "FRONTEND_HOST"
      )}/settings/organization/billing`,
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
    return this.stripe.customers.retrieve(stripeCustomerId, {
      expand: ["invoice_settings.default_payment_method"],
    }) as Promise<Stripe.Customer | Stripe.DeletedCustomer>;
  }

  async getPrice(id: string) {
    return this.stripe.prices.retrieve(id, {
      expand: ["product"],
    });
  }

  async getProduct(id: string) {
    return this.stripe.products.retrieve(id);
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

    return plans;
  }

  async updateCustomerDefaultPaymentMethod(
    customerId: string,
    paymentMethodId: string
  ) {
    await this.stripe.customers.update(customerId, {
      invoice_settings: {
        default_payment_method: paymentMethodId,
      },
    });
  }

  async updateSubscription(customerId: string, newPriceId: string) {
    // Retrieve the customer's subscriptions
    const subscriptions = await this.stripe.subscriptions.list({
      customer: customerId,
      expand: ["data.default_payment_method"],
      status: "active",
    });

    if (subscriptions.data.length === 0) {
      throw new BadRequestException(
        "No active subscriptions found for this customer."
      );
    }

    // Assuming there is only one subscription per customer
    const subscription = subscriptions.data[0];

    // Update the subscription to the new price
    await this.stripe.subscriptions.update(subscription.id, {
      items: [
        {
          id: subscription.items.data[0].id,
          price: newPriceId,
        },
      ],
      proration_behavior: "create_prorations",
    });
  }
}
