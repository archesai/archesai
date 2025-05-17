import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for customers.
 */
export class CustomersService {
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public async create(name: string, billingEmail: string) {
    return this.stripeService.stripe.customers.create({
      email: billingEmail,
      name
    })
  }

  public async findOne(id: string) {
    return this.stripeService.stripe.customers.retrieve(id, {
      expand: ['invoice_settings.default_payment_method']
    })
  }
}
