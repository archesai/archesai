import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import { InternalServerErrorException } from '@archesai/core'

import type { CallbacksService } from '#callbacks/callbacks.service'

export interface CallbacksControllerOptions {
  callbacksService: CallbacksService
}

export const callbacksController: FastifyPluginAsyncTypebox<
  CallbacksControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { callbacksService }) => {
  app.post(
    `/billing/stripe/callback`,
    {
      schema: {
        description: `Handles Stripe webhook callbacks`,
        hide: true,
        operationId: 'stripeCallback',
        response: {
          200: {
            description: 'OK'
          }
        },
        summary: `Stripe webhook callback`,
        tags: ['Billing - Callbacks']
      }
    },
    async (req) => {
      const stripeSignature = req.headers['stripe-signature']
      if (typeof stripeSignature !== 'string') {
        throw new InternalServerErrorException(
          'Missing stripe-signature header'
        )
      }
      await callbacksService.handle(stripeSignature, req.raw)
    }
  )
}
