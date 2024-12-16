import Stripe from 'stripe'

export class PlanEntity {
  /**
   * The currency of the plan
   * @example 'usd'
   */
  currency: string

  /**
   * The description of the plan
   * @example 'A plan for a small business'
   */
  description: null | string

  /**
   * The ID of the plan
   * @example 'prod_1234567890'
   */
  id: string

  /**
   * The metadata of the plan
   * @example { 'key': 'value' }
   */
  metadata: Record<string, string>

  /**
   * The name of the plan
   * @example 'Small Business Plan'
   */
  name: string

  /**
   * The ID of the price associated with the plan
   * @example 'price_1234567890'
   */
  priceId: string

  /**
   * The metadata of the price associated with the plan
   * @example { 'key': 'value' }
   */
  priceMetadata: Record<string, string>

  /**
   * The interval of the plan
   */
  recurring: null | Stripe.Price.Recurring

  /**
   * The amount in cents to be charged on the interval specified
   * @example 1000
   */
  unitAmount: number

  constructor(partial: Partial<PlanEntity>) {
    Object.assign(this, partial)
  }
}
