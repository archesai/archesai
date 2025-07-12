import type { StaticDecode } from '@sinclair/typebox'

import type { Controller, HttpInstance } from '@archesai/core'

import {
  createCollectionResponseSchema,
  createResourceObjectSchema,
  IS_CONTROLLER
} from '@archesai/core'
import { PLAN_ENTITY_KEY, PlanDtoSchema } from '@archesai/schemas'

import type { PlansService } from '#plans/plans.service'

const PlanResourceObjectSchema = createResourceObjectSchema(
  PlanDtoSchema,
  PLAN_ENTITY_KEY
)

const PlanCollectionResponseSchema = createCollectionResponseSchema(
  PlanResourceObjectSchema,
  PLAN_ENTITY_KEY
)

type PlanPaginatedResponse = StaticDecode<typeof PlanCollectionResponseSchema>

/**
 * Controller for Plans.
 */
export class PlansController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly plansService: PlansService

  constructor(plansService: PlansService) {
    this.plansService = plansService
  }

  public async findMany(): Promise<PlanPaginatedResponse> {
    const plans = await this.plansService.findAll()
    return {
      data: plans.data.map((plan) => {
        const { id, ...attributes } = plan
        return {
          attributes: attributes,
          id: id,
          links: {
            self: `/billing/${PLAN_ENTITY_KEY}/${plan.id}`
          },
          type: PLAN_ENTITY_KEY
        }
      })
    }
  }

  public registerRoutes(app: HttpInstance) {
    app.get(
      `/billing/${PLAN_ENTITY_KEY}`,
      {
        schema: {
          description: 'Get all plans',
          operationId: 'getPlans',
          response: {
            200: PlanCollectionResponseSchema
          },
          summary: 'Get all plans',
          tags: ['Billing']
        }
      },
      this.findMany.bind(this)
    )
  }
}
