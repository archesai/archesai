import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { Logger } from '@archesai/core'

import { PLAN_ENTITY_KEY, PlanDtoSchema, Type } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

import { PlansService } from '#plans/plans.service'

export interface PlansControllerOptions {
  logger: Logger
  stripeService: StripeService
}

export const plansController: FastifyPluginAsyncTypebox<
  PlansControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { logger, stripeService }) => {
  const plansService = new PlansService(stripeService, logger)
  app.get(
    `/billing/${PLAN_ENTITY_KEY}`,
    {
      schema: {
        description: 'Get all plans',
        operationId: 'getPlans',
        response: {
          200: Type.Array(PlanDtoSchema)
        },
        summary: 'Get all plans',
        tags: ['Billing']
      }
    },
    async () => {
      const plans = await plansService.findAll()
      return plans.data
    }
  )
}
