import type Stripe from 'stripe'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for customers.
 */
export class CustomersService {
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public async create(
    name: string,
    billingEmail: string
  ): Promise<Stripe.Response<Stripe.Customer>> {
    return this.stripeService.stripe.customers.create({
      email: billingEmail,
      name
    })
  }

  public async findOne(
    id: string
  ): Promise<Stripe.Response<Stripe.Customer | Stripe.DeletedCustomer>> {
    return this.stripeService.stripe.customers.retrieve(id, {
      expand: ['invoice_settings.default_payment_method']
    })
  }
}
