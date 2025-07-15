import type { Static, TLiteral, TObject, TUnion } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { MemberEntitySchema } from '#auth/members/entities/member.entity'

export const CreateMemberDtoSchema: TObject<{
  role: TUnion<TLiteral<'ADMIN' | 'USER'>[]>
}> = Type.Object({
  role: MemberEntitySchema.properties.role
})

export type CreateMemberDto = Static<typeof CreateMemberDtoSchema>
