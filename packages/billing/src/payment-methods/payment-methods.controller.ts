import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'

import { PaymentMethodEntitySchema, Type } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

import { CustomersService } from '#customers/customers.service'
import { PaymentMethodRepository } from '#payment-methods/payment-method.repository'
import { PaymentMethodsService } from '#payment-methods/payment-methods.service'

export interface PaymentMethodsControllerOptions {
  stripeService: StripeService
  websocketsService: WebsocketsService
}

export const paymentMethodsController: FastifyPluginAsyncTypebox<
  PaymentMethodsControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { stripeService, websocketsService }) => {
  const paymentMethodRepository = new PaymentMethodRepository(stripeService)
  const customersService = new CustomersService(stripeService)
  const paymentMethodsService = new PaymentMethodsService(
    customersService,
    paymentMethodRepository,
    websocketsService
  )
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
