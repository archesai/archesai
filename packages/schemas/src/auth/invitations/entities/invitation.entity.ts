import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
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
    organizationId: Type.String({
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

export type InvitationEntity = Static<typeof InvitationEntitySchema>

export const INVITATION_ENTITY_KEY = 'invitations'
