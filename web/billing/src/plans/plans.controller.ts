import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { Logger } from '@archesai/core'

import { PLAN_ENTITY_KEY, PlanDtoSchema } from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

import { PlansService } from '#plans/plans.service'

export interface PlansControllerOptions {
  logger: Logger
  stripeService: StripeService
}

export const plansController: FastifyPluginAsyncZod<
  PlansControllerOptions
> = async (app, { logger, stripeService }) => {
  const plansService = new PlansService(stripeService, logger)
  app.get(
    `/billing/${PLAN_ENTITY_KEY}`,
    {
      schema: {
        description: 'Get all plans',
        operationId: 'getPlans',
        response: {
          200: PlanDtoSchema.array()
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

  await Promise.resolve()
}
