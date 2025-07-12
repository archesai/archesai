import type { PlanDto } from '@archesai/schemas'

import { Logger } from '@archesai/core'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for Plans.
 */
export class PlansService {
  private readonly logger = new Logger(PlansService.name)
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public async findAll(): Promise<{
    count: number
    data: PlanDto[]
  }> {
    const products = await this.stripeService.stripe.products.list({
      active: true,
      expand: ['data.default_price']
    })

    const plans = products.data
      .map((product) => {
        const price = product.default_price
        if (!price) {
          return null
        }
        if (typeof price === 'string') {
          this.logger.warn('product without default price', { product })
          return null
        }
        return {
          createdAt: new Date(product.created).toISOString(),
          currency: price.currency,
          description: product.description,
          id: product.id,
          metadata: product.metadata,
          name: product.name,
          priceId: price.id,
          priceMetadata: price.metadata,
          recurring: price.recurring,
          slug: price.id,
          type: 'plan',
          unitAmount: price.unit_amount,
          updatedAt: new Date(product.updated).toISOString()
        } satisfies PlanDto
      })
      .filter((val) => val !== null)

    return {
      count: plans.length,
      data: plans
    }
  }
}
