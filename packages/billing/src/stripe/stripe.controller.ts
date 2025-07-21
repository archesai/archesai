import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import {
  CheckoutSessionDtoSchema,
  CreateCheckoutSessionDtoSchema,
  CreatePortalDtoSchema,
  PortalDtoSchema
} from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

export interface StripeControllerOptions {
  stripeService: StripeService
}

export const stripeController: FastifyPluginAsyncZod<
  StripeControllerOptions
> = async (app, { stripeService }) => {
  app.post(
    `/billing/portal`,
    {
      schema: {
        body: CreatePortalDtoSchema,
        description: 'Create a new portal',
        operationId: 'createPortal',
        response: {
          201: {
            description: 'The created portal',
            schema: PortalDtoSchema
          }
        },
        summary: 'Create a new portal',
        tags: ['Billing']
      }
    },
    (req) => {
      return stripeService.createPortal({
        organizationId: req.body.organizationId
      })
    }
  )

  app.post(
    `/billing/checkout-sessions`,
    {
      schema: {
        body: CreateCheckoutSessionDtoSchema,
        description: 'Create a checkout session',
        operationId: 'createCheckoutSession',
        response: {
          200: CheckoutSessionDtoSchema
        },
        summary: 'Create a checkout session',
        tags: ['Billing']
      }
    },
    (req) => {
      return stripeService.createCheckoutSession(
        '',
        {
          price: req.body.priceId,
          quantity: 1
        },
        false
      )
    }
  )

  await Promise.resolve()
}
