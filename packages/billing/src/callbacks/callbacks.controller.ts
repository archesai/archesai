import type { FastifyPluginCallbackTypebox } from '@fastify/type-provider-typebox'

import type { ConfigService, WebsocketsService } from '@archesai/core'

import { InternalServerErrorException } from '@archesai/core'

import type { StripeService } from '#stripe/stripe.service'

import { CallbacksService } from '#callbacks/callbacks.service'

export interface CallbacksControllerOptions {
  configService: ConfigService
  stripeService: StripeService
  websocketsService: WebsocketsService
}

export const callbacksController: FastifyPluginCallbackTypebox<
  CallbacksControllerOptions
> = (app, { configService, stripeService, websocketsService }, done) => {
  const callbacksService = new CallbacksService(
    configService,
    stripeService,
    websocketsService
  )
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

  done()
}
