import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { PlanTypes } from '#enums/role'

export const OrganizationEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    billingEmail: Type.String({
      description: 'The billing email to use for the organization'
    }),
    creator: Type.Optional(
      Type.String({
        description: 'The user who created the organization'
      })
    ),
    credits: Type.Number({
      description:
        'The number of credits you have remaining for this organization'
    }),
    customerId: Type.Optional(
      Type.String({
        description: 'The Stripe customer ID'
      })
    ),
    orgname: Type.String({
      description: 'The organization name'
    }),
    plan: Type.Union(
      PlanTypes.map((plan) => Type.Literal(plan)),
      { description: 'The plan that the organization is subscribed to' }
    )
  },
  {
    $id: 'OrganizationEntity',
    description: 'The organization entity',
    title: 'Organization Entity'
  }
)

export type OrganizationEntity = Static<typeof OrganizationEntitySchema>

export const ORGANIZATION_ENTITY_KEY = 'organizations'
