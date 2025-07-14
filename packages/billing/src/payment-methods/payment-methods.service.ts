import type { SearchQuery, WebsocketsService } from '@archesai/core'
import type { PaymentMethodEntity } from '@archesai/schemas'

import { BadRequestException, NotFoundException } from '@archesai/core'

import type { CustomersService } from '#customers/customers.service'
import type { PaymentMethodRepository } from '#payment-methods/payment-method.repository'

/**
 * Service for Payment Methods.
 */
export class PaymentMethodsService {
  private readonly customersService: CustomersService
  private readonly paymentMethodRepository: PaymentMethodRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    customersService: CustomersService,
    paymentMethodRepository: PaymentMethodRepository,
    websocketsService: WebsocketsService
  ) {
    this.customersService = customersService
    this.paymentMethodRepository = paymentMethodRepository
    this.websocketsService = websocketsService
  }

  public async delete(id: string): Promise<PaymentMethodEntity> {
    const paymentToDelete = await this.findOne(id)
    const paymentMethods = await this.findMany({
      filter: {
        customer: {
          equals: paymentToDelete.customer
        }
      }
    })
    if (paymentMethods.data.length <= 1) {
      throw new BadRequestException(
        'Cannot delete the last payment method. At least one payment method is required.'
      )
    }
    const customer = await this.customersService.findOne(
      paymentToDelete.customer
    )
    if (customer.deleted) {
      throw new NotFoundException('Customer has been deleted.')
    }

    // If the payment method being removed is the default, set another default
    if (customer.invoice_settings.default_payment_method === id) {
      const otherPaymentMethod = paymentMethods.data.find((pm) => pm.id !== id)
      if (otherPaymentMethod) {
        await this.updateCustomerId(customer.id, otherPaymentMethod.id)
      } else {
        throw new BadRequestException(
          'Cannot remove the last payment method. At least one payment method is required.'
        )
      }
    }

    const paymentMethod = await this.paymentMethodRepository.delete(
      paymentToDelete.id
    )

    this.websocketsService.broadcastEvent(paymentToDelete.customer, 'update', {
      queryKey: ['organizations', paymentToDelete.customer, 'billing']
    })

    return paymentMethod
  }

  public async findMany(query: SearchQuery<PaymentMethodEntity>): Promise<{
    count: number
    data: PaymentMethodEntity[]
  }> {
    return this.paymentMethodRepository.findMany(query)
  }

  public async findOne(id: string): Promise<PaymentMethodEntity> {
    return this.paymentMethodRepository.findOne(id)
  }

  public async updateCustomerId(id: string, customerId: string): Promise<void> {
    await this.paymentMethodRepository.update(id, {
      customer: customerId
    })
  }
}
