import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import { IdParamsSchema, UpdateSubscriptionDtoSchema } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

import { SubscriptionsService } from '#subscriptions/subscriptions.service'

export interface SubscriptionsControllerOptions {
  stripeService: StripeService
}

export const subscriptionsController: FastifyPluginAsyncZod<
  SubscriptionsControllerOptions
> = async (app, { stripeService }) => {
  const subscriptionsService = new SubscriptionsService(stripeService)
  app.delete(
    `/billing/subscriptions/:id`,
    {
      schema: {
        description: 'Cancel a subscription',
        operationId: 'cancelSubscription',
        params: IdParamsSchema,
        summary: 'Cancel a subscription',
        tags: ['Billing']
      }
    },
    async (req) => {
      await subscriptionsService.cancel(req.params.id)
    }
  )

  app.patch(
    `/billing/subscriptions/:id`,
    {
      schema: {
        body: UpdateSubscriptionDtoSchema,
        description: 'Update a subscription',
        operationId: 'updateSubscription',
        params: IdParamsSchema,
        summary: 'Update a subscription',
        tags: ['Billing']
      }
    },
    (req) => {
      return subscriptionsService.update(req.params.id, req.body.planId)
    }
  )

  await Promise.resolve()
}
