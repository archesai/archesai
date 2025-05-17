import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { RoleType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const InvitationEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    accepted: Type.Boolean({
      description: 'Whether the invite was accepted'
    }),
    email: Type.String({
      description: 'The email of the invitated user'
    }),
    orgname: Type.String({
      description: 'The name of the organization the token belongs to'
    }),
    role: Type.Union(
      RoleTypes.map((role) => Type.Literal(role)), // Using literals instead of enums
      { description: 'The role of the invitation' }
    )
  },
  {
    $id: 'InvitationEntity',
    description: 'The invitation entity',
    title: 'Invitation Entity'
  }
)

export class InvitationEntity
  extends BaseEntity
  implements Static<typeof InvitationEntitySchema>
{
  public accepted: boolean
  public email: string
  public orgname: string
  public role: RoleType
  public type = INVITATION_ENTITY_KEY

  constructor(props: InvitationEntity) {
    super(props)
    this.accepted = props.accepted
    this.email = props.email
    this.orgname = props.orgname
    this.role = props.role
  }
}

export const INVITATION_ENTITY_KEY = 'invitations'
