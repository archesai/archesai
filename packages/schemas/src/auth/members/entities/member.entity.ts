import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const MemberEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    invitationId: Type.String({ description: 'The invitation id' }),
    organizationId: Type.String({ description: 'The organization name' }),
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

export type MemberEntity = Static<typeof MemberEntitySchema>

export const MEMBER_ENTITY_KEY = 'members'
