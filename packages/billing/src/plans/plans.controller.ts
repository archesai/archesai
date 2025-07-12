import type { Controller, HttpInstance } from '@archesai/core'

import { IS_CONTROLLER } from '@archesai/core'
import { PLAN_ENTITY_KEY, PlanDtoSchema, Type } from '@archesai/schemas'

import type { PlansService } from '#plans/plans.service'

/**
 * Controller for Plans.
 */
export class PlansController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly plansService: PlansService

  constructor(plansService: PlansService) {
    this.plansService = plansService
  }

  public registerRoutes(app: HttpInstance) {
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
        const { data: plans } = await this.plansService.findAll()
        return plans
      }
    )
  }
}
