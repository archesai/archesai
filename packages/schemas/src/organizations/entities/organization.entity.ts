import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { PlanType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
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

export class OrganizationEntity
  extends BaseEntity
  implements Static<typeof OrganizationEntitySchema>
{
  public billingEmail: string
  public credits: number
  public customer?: string
  public orgname: string
  public plan: PlanType
  public type = ORGANIZATION_ENTITY_KEY

  constructor(props: OrganizationEntity) {
    super(props)
    this.billingEmail = props.billingEmail
    this.credits = props.credits
    if (props.customer) {
      this.customer = props.customer
    }
    this.orgname = props.orgname
    this.plan = props.plan
  }
}

export const ORGANIZATION_ENTITY_KEY = 'organizations'
