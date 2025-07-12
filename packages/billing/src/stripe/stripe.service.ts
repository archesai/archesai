import Stripe from 'stripe'

import type { ConfigService } from '@archesai/core'
import type { CreatePortalDto, PortalDto } from '@archesai/schemas'

import { InternalServerErrorException } from '@archesai/core'

/**
 * Service for communicating with Stripe.
 */
export class StripeService {
  public readonly configService: ConfigService
  public readonly stripe: Stripe

  constructor(configService: ConfigService) {
    this.stripe = new Stripe(configService.get('billing.stripe.token'), {
      apiVersion: '2025-06-30.basil'
    })
    this.configService = configService
  }

  public async createCheckoutSession(
    customerId: string,
    lineItem: { price: string; quantity: number },
    isOneTime: boolean
  ) {
    const session = await this.stripe.checkout.sessions.create({
      ...(isOneTime ? { allow_promotion_codes: true } : {}),
      cancel_url: `${this.configService.get('platform.host')}/organization/billing`,
      customer: customerId,
      ...(isOneTime ?
        {
          invoice_creation: { enabled: true }
        }
      : {}),

      line_items: [
        {
          ...lineItem,
          ...(isOneTime ?
            { adjustable_quantity: { enabled: true, minimum: 1 } }
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

  public async createPortal(
    createPortalRequest: CreatePortalDto
  ): Promise<PortalDto> {
    return this.stripe.billingPortal.sessions.create({
      customer: createPortalRequest.organizationId,
      return_url: `${this.configService.get('platform.host')}/organization/billing`
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
