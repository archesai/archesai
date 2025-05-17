import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema, PlanTypes } from '@archesai/domain'

export const PLAN_ENTITY_KEY = 'plans'

export const PlanResourceSchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    currency: Type.String({ description: 'The currency of the plan' }),
    description: Type.Optional(
      Type.Union([Type.String(), Type.Null()], {
        description: 'The description of the plan'
      })
    ),
    id: Type.String({ description: 'The ID of the plan' }),
    metadata: Type.Object({
      key: Type.Optional(
        Type.Union(
          PlanTypes.map((plan) => Type.Literal(plan)),
          { description: 'The key of the metadata' }
        )
      )
    }),
    name: Type.String({ description: 'The name of the plan' }),
    priceId: Type.String({
      description: 'The ID of the price associated with the plan'
    }),
    priceMetadata: Type.Record(Type.String(), Type.String(), {
      description: 'The metadata of the price associated with the plan'
    }),
    recurring: Type.Optional(
      Type.Union(
        [
          Type.Object({
            interval: Type.String(),
            interval_count: Type.Number(),
            trial_period_days: Type.Optional(
              Type.Union([Type.Number(), Type.Null()])
            )
          }),
          Type.Null()
        ],
        { description: 'The interval of the plan' }
      )
    ),

    unitAmount: Type.Optional(
      Type.Union([Type.Number(), Type.Null()], {
        description:
          'The amount in cents to be charged on the interval specified'
      })
    )
  },
  {
    $id: 'PlanResource',
    description: 'The plan resource',
    title: 'Plan Resource'
  }
)

export type PlanResource = Static<typeof PlanResourceSchema>
