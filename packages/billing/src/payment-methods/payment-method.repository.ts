import type Stripe from 'stripe'

import type { BaseRepository, SearchQuery } from '@archesai/core'

import { NotFoundException } from '@archesai/core'
import { generateSlug, PaymentMethodEntity } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Repository for Payment Methods.
 */
export class PaymentMethodRepository
  implements
    Pick<
      BaseRepository<
        PaymentMethodEntity,
        PaymentMethodEntity,
        Stripe.PaymentMethod
      >,
      | 'create'
      | 'createMany'
      | 'delete'
      | 'deleteMany'
      | 'findMany'
      | 'findOne'
      | 'update'
      | 'updateMany'
    >
{
  protected readonly primaryKey = 'payment-method'
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public async create(
    _value: PaymentMethodEntity
  ): Promise<PaymentMethodEntity> {
    const pm = await this.stripeService.stripe.paymentMethods.create({
      type: 'card'
    })
    return this.toEntity(pm)
  }

  public async createMany(
    values: PaymentMethodEntity[]
  ): Promise<{ count: number; data: PaymentMethodEntity[] }> {
    const created: PaymentMethodEntity[] = []
    for (const val of values) {
      const pm = await this.create(val)
      created.push(pm)
    }
    return {
      count: 0,
      data: created
    }
  }

  public async delete(id: string): Promise<PaymentMethodEntity> {
    const detached = await this.stripeService.stripe.paymentMethods.detach(id)
    return this.toEntity(detached)
  }

  public async deleteMany(
    query: SearchQuery<PaymentMethodEntity>
  ): Promise<{ count: number; data: PaymentMethodEntity[] }> {
    const found = await this.findMany(query)
    for (const pm of found.data) {
      await this.delete(pm.id)
    }
    return { count: found.count, data: found.data }
  }

  public async findFirst(
    query: SearchQuery<PaymentMethodEntity>
  ): Promise<PaymentMethodEntity> {
    const { data } = await this.findMany(query)
    const first = data[0]
    if (!first) {
      throw new NotFoundException('No payment method found')
    }
    return first
  }

  public async findMany(
    query: SearchQuery<PaymentMethodEntity>
  ): Promise<{ count: number; data: PaymentMethodEntity[] }> {
    const customerId = query.filter?.customer?.equals
    const results = await this.stripeService.stripe.paymentMethods.list({
      customer: typeof customerId === 'string' ? customerId : '',
      type: 'card'
    })

    return {
      count: results.data.length,
      data: results.data.map((pm) => this.toEntity(pm))
    }
  }

  public async findOne(id: string): Promise<PaymentMethodEntity> {
    const pm = await this.stripeService.stripe.paymentMethods.retrieve(id, {
      expand: ['customer']
    })

    if (pm.id !== id) {
      throw new NotFoundException(`Payment method ${id} not found`)
    }
    return this.toEntity(pm)
  }

  public async update(
    id: string,
    _: Partial<PaymentMethodEntity>
  ): Promise<PaymentMethodEntity> {
    // Example: update "billing_details" or "metadata" in Stripe
    const updated = await this.stripeService.stripe.paymentMethods.update(id, {
      // e.g., billing_details: data.billing_details
    })

    return this.toEntity(updated)
  }

  public async updateMany(
    value: Partial<PaymentMethodEntity>,
    query: SearchQuery<PaymentMethodEntity>
  ): Promise<{ count: number; data: PaymentMethodEntity[] }> {
    const found = await this.findMany(query)
    const updated: PaymentMethodEntity[] = []

    for (const pm of found.data) {
      const u = await this.update(pm.id, value)
      updated.push(u)
    }

    return {
      count: updated.length,
      data: updated
    }
  }

  protected toEntity(model: Stripe.PaymentMethod): PaymentMethodEntity {
    return new PaymentMethodEntity({
      ...model,
      createdAt: model.created.toString(),
      customer: model.customer as string,
      name: model.billing_details.name ?? '',
      slug: generateSlug(model.id),
      stripeId: model.id,
      type: 'payment-method',
      updatedAt: model.created.toString()
    })
  }
}
