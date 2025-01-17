import {
  BadRequestException,
  Controller,
  Delete,
  Get,
  NotFoundException,
  Param
} from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'
import { Stripe } from 'stripe'
import { OrganizationsService } from '@/src/organizations/organizations.service'
import { WebsocketsService } from '@/src/websockets/websockets.service'
import { BillingService } from '@/src/billing/billing.service'
import { PaymentMethodEntity } from '@/src/billing/entities/payment-method.entity'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'
import { RoleTypeEnum } from '@/src/members/entities/member.entity'

@ApiTags('Billing - Payment Methods')
@Authenticated([RoleTypeEnum.ADMIN])
@Controller('/organizations/:orgname/payment-methods')
export class PaymentMethodsController {
  constructor(
    private billingService: BillingService,
    private organizationsService: OrganizationsService,
    private websocketsService: WebsocketsService
  ) {}

  /**
   * List payment methods
   * @remarks This endpoint will return a list of payment methods for an organization
   */
  @Get()
  async findAll(@Param('orgname') orgname: string) {
    const organization = await this.organizationsService.findOne(orgname)
    const paymentMethods = await this.billingService.listPaymentMethods(
      organization.stripeCustomerId
    )
    return paymentMethods.data.map((pm) => new PaymentMethodEntity(pm))
  }

  /**
   * Remove payment method
   * @remarks This endpoint will remove a payment method from an organization
   * @throws {404} NotFoundException
   */
  @Delete(':paymentMethodId')
  async remove(
    @Param('orgname') orgname: string,
    @Param('paymentMethodId') paymentMethodId: string
  ) {
    const organization = await this.organizationsService.findOne(orgname)
    const paymentMethods = await this.billingService.listPaymentMethods(
      organization.stripeCustomerId
    )
    const paymentMethod = paymentMethods.data.find(
      (pm) => pm.id === paymentMethodId
    )
    if (!paymentMethod) {
      throw new NotFoundException('Payment method not found')
    }

    // Check if there is more than one payment method
    if (paymentMethods.data.length <= 1) {
      throw new BadRequestException(
        'Cannot remove the last payment method. At least one payment method is required.'
      )
    }

    // Retrieve the customer to check default payment method
    const customer = await this.billingService.getCustomer(
      organization.stripeCustomerId
    )

    // Type guard to ensure customer is not deleted
    if (customer.deleted) {
      throw new NotFoundException('Customer has been deleted.')
    }

    const c = customer as Stripe.Customer
    let defaultPaymentMethodId
    if (c.invoice_settings && c.invoice_settings.default_payment_method) {
      if (typeof c.invoice_settings.default_payment_method === 'string') {
        defaultPaymentMethodId = c.invoice_settings.default_payment_method
      } else {
        defaultPaymentMethodId = c.invoice_settings.default_payment_method.id
      }
    }

    // If the payment method being removed is the default, set another as default
    if (defaultPaymentMethodId === paymentMethodId) {
      // Set another payment method as default
      const otherPaymentMethod = paymentMethods.data.find(
        (pm) => pm.id !== paymentMethodId
      )
      if (otherPaymentMethod) {
        await this.billingService.updateCustomerDefaultPaymentMethod(
          organization.stripeCustomerId,
          otherPaymentMethod.id
        )
      } else {
        throw new BadRequestException(
          'Cannot remove the last payment method. At least one payment method is required.'
        )
      }

      this.websocketsService.socket?.to(orgname).emit('update', {
        queryKey: ['organizations', orgname, 'billing']
      })
    }

    await this.billingService.detachPaymentMethod(paymentMethodId)
  }
}
