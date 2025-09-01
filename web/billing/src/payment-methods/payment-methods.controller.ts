import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'

import { IdParamsSchema, PaymentMethodEntitySchema } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

import { CustomersService } from '#customers/customers.service'
import { PaymentMethodRepository } from '#payment-methods/payment-method.repository'
import { PaymentMethodsService } from '#payment-methods/payment-methods.service'

export interface PaymentMethodsControllerOptions {
  stripeService: StripeService
  websocketsService: WebsocketsService
}

export const paymentMethodsController: FastifyPluginAsyncZod<
  PaymentMethodsControllerOptions
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
        params: IdParamsSchema,
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
        params: IdParamsSchema,
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

  await Promise.resolve()
}
