import type {
  Static,
  TLiteral,
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

export const OrganizationEntitySchema: TObject<{
  billingEmail: TString
  createdAt: TString
  credits: TNumber
  id: TString
  logo: TOptional<TString>
  metadata: TOptional<TRecord<TString, TString>>
  name: TString
  plan: TUnion<
    TLiteral<'BASIC' | 'FREE' | 'PREMIUM' | 'STANDARD' | 'UNLIMITED'>[]
  >
  slug: TString
  stripeCustomerId: TOptional<TString>
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    billingEmail: Type.String({
      description: 'The billing email to use for the organization'
    }),
    credits: Type.Number({
      description:
        'The number of credits you have remaining for this organization'
    }),
    logo: Type.Optional(
      Type.String({
        description: 'The URL of the organization logo'
      })
    ),
    metadata: Type.Optional(
      Type.Record(Type.String(), Type.String(), {
        description: 'The metadata for the organization, used for custom data'
      })
    ),
    name: Type.String({
      description: 'The name of the organization'
    }),
    plan: Type.Union(
      PlanTypes.map((plan) => Type.Literal(plan)),
      { description: 'The plan that the organization is subscribed to' }
    ),
    slug: Type.String({
      description: 'The unique slug for the organization, used in URLs'
    }),
    stripeCustomerId: Type.Optional(
      Type.String({
        description: 'The Stripe customer ID'
      })
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
