import type { ConfigService } from '@archesai/core'

import { InternalServerErrorException } from '@archesai/core'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for checkout sessions.
 */
export class CheckoutSessionsService {
  private readonly configService: ConfigService
  private readonly stripeService: StripeService

  constructor(configService: ConfigService, stripeService: StripeService) {
    this.configService = configService
    this.stripeService = stripeService
  }

  public async create(
    customerId: string,
    lineItem: { price: string; quantity: number },
    isOneTime: boolean
  ) {
    const session = await this.stripeService.stripe.checkout.sessions.create({
      ...(isOneTime ? { allow_promotion_codes: true } : {}),
      cancel_url: `${this.configService.get('platform.host')}/organization/billing`,
      customer: customerId,
      ...(isOneTime
        ? {
            invoice_creation: { enabled: true }
          }
        : {}),

      line_items: [
        {
          ...lineItem,
          ...(isOneTime
            ? { adjustable_quantity: { enabled: true, minimum: 1 } }
            : {})
        }
      ],
      mode: isOneTime ? 'payment' : 'subscription',
      payment_method_types: ['card'],
      success_url: `${this.configService.get('platform.host')}/organization/billing`
    })

    if (!session.url) {
      throw new InternalServerErrorException(
        'Failed to create checkout session.'
      )
    }

    return { url: session.url }
  }
}
