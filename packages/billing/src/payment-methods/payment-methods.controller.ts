import { Type } from '@sinclair/typebox'

import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type { PaymentMethodEntity } from '@archesai/domain'

import { IS_CONTROLLER } from '@archesai/core'
import { PaymentMethodEntitySchema } from '@archesai/domain'

import type { PaymentMethodsService } from '#payment-methods/payment-methods.service'

export class PaymentMethodsController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly paymentMethodsService: PaymentMethodsService

  constructor(paymentMethodsService: PaymentMethodsService) {
    this.paymentMethodsService = paymentMethodsService
  }

  public async delete(request: ArchesApiRequest & { params: { id: string } }) {
    return this.paymentMethodsService.delete(request.params.id)
  }

  public async findMany(request: ArchesApiRequest) {
    return this.paymentMethodsService.findMany({
      filter: {
        customer: {
          equals: request.user!.id
        }
      }
    })
  }

  public async findOne(
    request: ArchesApiRequest & { params: { id: string } }
  ): Promise<PaymentMethodEntity> {
    return this.paymentMethodsService.findOne(request.params.id)
  }

  public registerRoutes(app: HttpInstance) {
    app.delete(
      `/billing/payment-methods/:id`,
      {
        schema: {
          description: 'Delete a payment method',
          operationId: 'deletePaymentMethod',
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: PaymentMethodEntitySchema
          },
          summary: 'Delete a payment method',
          tags: ['Billing - Payment Methods']
        }
      },
      this.delete.bind(this)
    )

    app.get(
      `/billing/payment-methods`,
      {
        schema: {
          description: 'Get all payment methods',
          operationId: 'findManyPaymentMethods',
          response: {
            200: {
              description: 'The payment methods of the customer'
            }
          },
          summary: 'Get all payment methods',
          tags: ['Billing - Payment Methods']
        }
      },
      this.findMany.bind(this)
    )

    app.get(
      `/billing/payment-methods/:id`,
      {
        schema: {
          description: 'Get a payment method',
          operationId: 'findOnePaymentMethod',
          params: Type.Object({
            id: Type.String()
          }),
          response: {
            200: PaymentMethodEntitySchema
          },
          summary: 'Get a payment method',
          tags: ['Billing - Payment Methods']
        }
      },
      (request) => {
        return this.paymentMethodsService.findOne(request.params.id)
      }
    )
  }
}
