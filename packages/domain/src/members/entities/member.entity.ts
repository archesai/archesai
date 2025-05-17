import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { RoleType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const MemberEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    invitationId: Type.String({ description: 'The invitation id' }),
    orgname: Type.String({ description: 'The organization name' }),
    role: Type.Union(
      RoleTypes.map((role) => Type.Literal(role)),
      { description: 'The role of the member' }
    ),
    userId: Type.String({ description: 'The user id' })
  },
  {
    $id: 'MemberEntity',
    description: 'The member entity',
    title: 'Member Entity'
  }
)

export class MemberEntity
  extends BaseEntity
  implements Static<typeof MemberEntitySchema>
{
  public invitationId: string
  public orgname: string
  public role: RoleType
  public type = MEMBER_ENTITY_KEY
  public userId: string

  constructor(props: MemberEntity) {
    super(props)
    this.invitationId = props.invitationId
    this.orgname = props.orgname
    this.role = props.role
    this.userId = props.userId
  }
}

export const MEMBER_ENTITY_KEY = 'members'
