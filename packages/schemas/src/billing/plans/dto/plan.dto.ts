import type {
  Static,
  TLiteral,
  TNull,
  TNumber,
  TObject,
  TOptional,
  TRecord,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { PlanTypes } from '#enums/role'

export const PLAN_ENTITY_KEY = 'plans'

export const PlanDtoSchema: TObject<{
  createdAt: TString
  currency: TString
  description: TOptional<TUnion<[TString, TNull]>>
  id: TString
  metadata: TObject<{
    key: TOptional<
      TUnion<
        TLiteral<'BASIC' | 'FREE' | 'PREMIUM' | 'STANDARD' | 'UNLIMITED'>[]
      >
    >
  }>
  name: TString
  priceId: TString
  priceMetadata: TRecord<TString, TString>
  recurring: TOptional<
    TUnion<
      [
        TObject<{
          interval: TString
          interval_count: TNumber
          trial_period_days: TOptional<TUnion<[TNumber, TNull]>>
        }>,
        TNull
      ]
    >
  >
  unitAmount: TOptional<TUnion<[TNumber, TNull]>>
  updatedAt: TString
}> = Type.Object(
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
    $id: 'PlanDto',
    description: 'The plan resource',
    title: 'Plan Resource'
  }
)

export type PlanDto = Static<typeof PlanDtoSchema>
