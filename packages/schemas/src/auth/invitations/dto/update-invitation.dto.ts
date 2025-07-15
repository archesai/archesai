import type {
  Static,
  TLiteral,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateInvitationDtoSchema } from '#auth/invitations/dto/create-invitation.dto'

export const UpdateInvitationDtoSchema: TObject<{
  email: TOptional<TString>
  role: TOptional<TUnion<TLiteral<'ADMIN' | 'USER'>[]>>
}> = Type.Partial(CreateInvitationDtoSchema)

export type UpdateInvitationDto = Static<typeof UpdateInvitationDtoSchema>
