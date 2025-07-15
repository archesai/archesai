import type {
  Static,
  TLiteral,
  TObject,
  TOptional,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateMemberDtoSchema } from '#auth/members/dto/create-member.dto'

export const UpdateMemberDtoSchema: TObject<{
  role: TOptional<TUnion<TLiteral<'ADMIN' | 'USER'>[]>>
}> = Type.Partial(CreateMemberDtoSchema)

export type UpdateMemberDto = Static<typeof UpdateMemberDtoSchema>
