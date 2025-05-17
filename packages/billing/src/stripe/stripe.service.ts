import Stripe from 'stripe'

import type { ConfigService } from '@archesai/core'

/**
 * Service for communicating with Stripe.
 */
export class StripeService {
  public readonly stripe: Stripe

  constructor(configService: ConfigService) {
    this.stripe = new Stripe(configService.get('billing.stripe.token'), {
      apiVersion: '2025-03-31.basil'
    })
  }

  public async createSetupIntent(customerId: string) {
    return this.stripe.setupIntents.create({
      customer: customerId,
      payment_method_types: ['card']
    })
  }

  public async getPrice(id: string) {
    return this.stripe.prices.retrieve(id, {
      expand: ['product']
    })
  }

  public async getProduct(id: string) {
    return this.stripe.products.retrieve(id)
  }
}
