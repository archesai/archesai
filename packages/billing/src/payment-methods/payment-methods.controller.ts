import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import { PaymentMethodEntitySchema, Type } from '@archesai/schemas'

import type { PaymentMethodsService } from '#payment-methods/payment-methods.service'

export interface PaymentMethodsControllerOptions {
  paymentMethodsService: PaymentMethodsService
}

export const paymentMethodsController: FastifyPluginAsyncTypebox<
  PaymentMethodsControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { paymentMethodsService }) => {
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
        tags: ['Billing']
      }
    },
    (req) => {
      return paymentMethodsService.delete(req.params.id)
    }
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
        tags: ['Billing']
      }
    },
    async (_req) => {
      return paymentMethodsService.findMany({
        // filter: {
        //   customer: {
        //     equals: req.user!.id
        //   }
        // }
      })
    }
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
        tags: ['Billing']
      }
    },
    async (req) => {
      return paymentMethodsService.findOne(req.params.id)
    }
  )
}
