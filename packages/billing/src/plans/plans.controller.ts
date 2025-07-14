import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import { PLAN_ENTITY_KEY, PlanDtoSchema, Type } from '@archesai/schemas'

import type { PlansService } from '#plans/plans.service'

export interface PlansControllerOptions {
  plansService: PlansService
}

export const plansController: FastifyPluginAsyncTypebox<
  PlansControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { plansService }) => {
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
