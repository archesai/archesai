import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for managing subscriptions.
 */
export class SubscriptionsService {
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public async cancel(id: string) {
    const subscription =
      await this.stripeService.stripe.subscriptions.retrieve(id)
    if (subscription.status === 'canceled') {
      return
    }

    await this.stripeService.stripe.subscriptions.cancel(id)
  }

  public async update(id: string, priceId: string) {
    const subscription =
      await this.stripeService.stripe.subscriptions.retrieve(id)

    return this.stripeService.stripe.subscriptions.update(subscription.id, {
      items: [
        {
          id,
          price: priceId
        }
      ],
      proration_behavior: 'create_prorations'
    })
  }
}
