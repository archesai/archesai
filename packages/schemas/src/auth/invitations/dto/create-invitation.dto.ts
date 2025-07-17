import type {
  Static,
  TLiteral,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { InvitationEntitySchema } from '#auth/invitations/entities/invitation.entity'

export const CreateInvitationDtoSchema: TObject<{
  email: TString
  role: TUnion<TLiteral<'admin' | 'member' | 'owner'>[]>
}> = Type.Object({
  email: InvitationEntitySchema.properties.email,
  role: InvitationEntitySchema.properties.role
})

export type CreateInvitationDto = Static<typeof CreateInvitationDtoSchema>
