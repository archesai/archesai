import type {
  Static,
  TLiteral,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const INVITATION_ENTITY_KEY = 'invitations'

export const InvitationEntitySchema: TObject<{
  createdAt: TString
  email: TString
  expiresAt: TString
  id: TString
  inviterId: TString
  organizationId: TString
  role: TUnion<TLiteral<'admin' | 'member' | 'owner'>[]>
  status: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    email: Type.String({
      description: 'The email of the invitated user'
    }),
    expiresAt: Type.String({
      description: 'The date and time when the invitation expires'
    }),
    inviterId: Type.String({
      description: 'The user id of the inviter'
    }),
    organizationId: Type.String({
      description: 'The name of the organization the token belongs to'
    }),
    role: Type.Union(
      RoleTypes.map((role) => Type.Literal(role)), // Using literals instead of enums
      { description: 'The role of the invitation' }
    ),
    status: Type.String({
      description:
        'The status of the invitation, e.g., pending, accepted, declined'
    })
  },
  {
    $id: 'InvitationEntity',
    description: 'The invitation entity',
    title: 'Invitation Entity'
  }
)

export type InvitationEntity = Static<typeof InvitationEntitySchema>
