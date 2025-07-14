import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import { Type, UpdateSubscriptionDtoSchema } from '@archesai/schemas'

import type { SubscriptionsService } from '#subscriptions/subscriptions.service'

export interface SubscriptionsControllerOptions {
  subscriptionsService: SubscriptionsService
}

export const subscriptionsController: FastifyPluginAsyncTypebox<
  SubscriptionsControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { subscriptionsService }) => {
  app.delete(
    `/billing/subscriptions/:id`,
    {
      schema: {
        description: 'Cancel a subscription',
        operationId: 'cancelSubscription',
        params: Type.Object({
          id: Type.String()
        }),
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
        params: Type.Object({
          id: Type.String()
        }),
        summary: 'Update a subscription',
        tags: ['Billing']
      }
    },
    (req) => {
      return subscriptionsService.update(req.params.id, req.body.planId)
    }
  )
}
