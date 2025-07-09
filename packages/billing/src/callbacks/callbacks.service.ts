import type Stripe from 'stripe'

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'

import type {
  ArchesApiRequest,
  ConfigService,
  EventBus,
  WebsocketsService
} from '@archesai/core'
import type {
  OrganizationCustomerSubscriptionUpdatedEvent,
  PlanType
} from '@archesai/schemas'

import { InternalServerErrorException } from '@archesai/core'
import { PlanTypes } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for callbacks.
 */
export class CallbacksService {
  private readonly configService: ConfigService
  private readonly eventBus: EventBus
  private readonly stripeService: StripeService
  private readonly websocketsService: WebsocketsService

  constructor(
    configService: ConfigService,
    eventBus: EventBus,
    stripeService: StripeService,
    websocketsService: WebsocketsService
  ) {
    this.configService = configService
    this.eventBus = eventBus
    this.stripeService = stripeService
    this.websocketsService = websocketsService
  }

  public async handle(signature: string, req: ArchesApiRequest) {
    if (!signature) {
      throw new InternalServerErrorException('Missing stripe-signature header')
    }

    const event = this.constructEventFromPayload(
      signature,
      Buffer.from(await req.raw.toArray())
    )

    const eventObj = event.data.object
    if (
      !('customer' in eventObj) ||
      !('metadata' in eventObj) ||
      !eventObj.metadata ||
      !('orgname' in eventObj.metadata)
    ) {
      throw new InternalServerErrorException('Invalid event object')
    }

    if (
      typeof eventObj.customer !== 'string' ||
      typeof eventObj.metadata.orgname !== 'string'
    ) {
      throw new InternalServerErrorException('Invalid event object')
    }
    const customer = eventObj.customer
    const orgname = eventObj.metadata.orgname

    if (event.type == 'invoice.paid') {
      const data = event.data.object
      if (data.amount_paid > 0) {
        const orgname = event.data.object.metadata?.orgname
        if (typeof orgname !== 'string') {
          throw new InternalServerErrorException('Invalid orgname')
        }
        for (const lineItem of data.lines.data) {
          const priceId = lineItem.id
          if (!priceId) {
            continue
          }
          const price = await this.stripeService.getPrice(priceId)
          const product = price.product
          if (
            typeof product === 'string' ||
            product.deleted ||
            !product.metadata.credits
          ) {
            throw new InternalServerErrorException('Invalid product')
          }
          const credits = product.metadata.credits
          const quantity = lineItem.quantity ?? 1

          this.eventBus.emit('organization.customer.subscription.updated', {
            credits: Number(credits) * quantity,
            customer: customer,
            orgname: orgname
          } satisfies OrganizationCustomerSubscriptionUpdatedEvent)

          this.websocketsService.broadcastEvent(orgname, 'update', {
            queryKey: ['organizations', orgname]
          })
        }
      }
    }

    if (
      event.type == 'customer.subscription.created' ||
      event.type == 'customer.subscription.updated' ||
      event.type == 'customer.subscription.deleted'
    ) {
      const data = event.data.object satisfies Stripe.Subscription
      const priceId = Value.Parse(Type.String(), data.items.data[0]?.price.id)
      const price = await this.stripeService.getPrice(priceId)
      if (typeof price.product === 'string' || price.product.deleted) {
        throw new InternalServerErrorException('Invalid price')
      }
      const planType = price.product.metadata.key
      if (!planType || !PlanTypes.includes(planType as PlanType)) {
        throw new InternalServerErrorException('Invalid plan type')
      }

      this.eventBus.emit('organization.customer.subscription.updated', {
        customer: customer,
        orgname: orgname,
        planType: event.data.object.status == 'active' ? planType : 'FREE'
      } satisfies OrganizationCustomerSubscriptionUpdatedEvent)

      this.websocketsService.broadcastEvent(orgname, 'update', {
        queryKey: ['organizations', orgname]
      })
    }
  }

  private constructEventFromPayload(signature: string, payload: Buffer) {
    const webhookSecret = this.configService.get('billing.stripe.whsec')
    return this.stripeService.stripe.webhooks.constructEvent(
      payload,
      signature,
      webhookSecret
    )
  }
}
