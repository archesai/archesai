import {
  BadRequestException,
  Controller,
  Headers,
  Post,
  RawBodyRequest,
  Req
} from '@nestjs/common'
import { ApiExcludeController } from '@nestjs/swagger'
import { PlanType } from '@prisma/client'
import { Stripe } from 'stripe'

import { OrganizationsService } from '@/src/organizations/organizations.service'
import { WebsocketsService } from '@/src/websockets/websockets.service'
import { BillingService } from '@/src/billing/billing.service'

@ApiExcludeController()
@Controller('billing/stripe')
export class StripeController {
  constructor(
    private billingService: BillingService,
    private organizationsService: OrganizationsService,
    private websocketsService: WebsocketsService
  ) {}

  @Post('callback')
  async stripeCallback(
    @Headers('stripe-signature') signature: string,
    @Req() req: RawBodyRequest<Request>
  ) {
    if (!signature) {
      throw new BadRequestException('Missing stripe-signature header')
    }

    if (!req.rawBody) {
      throw new BadRequestException('Missing raw body')
    }

    const event = await this.billingService.constructEventFromPayload(
      signature,
      req.rawBody
    )

    if (event.type == 'invoice.paid') {
      const data = event.data.object as Stripe.Invoice
      if (data.amount_paid > 0) {
        const customerId = data.customer as string
        const organization =
          await this.organizationsService.findByStripeCustomerId(customerId)
        for (const lineItem of data.lines.data) {
          const priceId = lineItem?.price?.id
          if (!priceId) {
            continue
          }
          const price = await this.billingService.getPrice(priceId)
          const product = price.product as Stripe.Product
          const credits = product.metadata['credits']
          const quantity = lineItem.quantity || 1
          await this.organizationsService.addOrRemoveCredits(
            organization.orgname,
            Number(credits) * quantity
          )
          this.websocketsService.socket
            ?.to(organization.orgname)
            .emit('update', {
              queryKey: ['organizations', organization.orgname]
            })
        }
      }
    }

    if (
      event.type == 'customer.subscription.created' ||
      event.type == 'customer.subscription.updated' ||
      event.type == 'customer.subscription.deleted'
    ) {
      const data = event.data.object as Stripe.Subscription
      const customerId = data.customer as string
      const organization =
        await this.organizationsService.findByStripeCustomerId(customerId)

      const priceId = data.items.data[0].price.id
      const price = await this.billingService.getPrice(priceId)
      const product = price.product as Stripe.Product
      const planType = product.metadata['key'] as PlanType

      if (data.status == 'active') {
        await this.organizationsService.setPlan(organization.orgname, planType)
      } else {
        await this.organizationsService.setPlan(organization.orgname, 'FREE')
      }
      this.websocketsService.socket?.to(organization.orgname).emit('update', {
        queryKey: ['organizations', organization.orgname]
      })
    }
  }
}
