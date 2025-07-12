import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { MemberEntitySchema } from '#members/entities/member.entity'

export const CreateMemberDtoSchema = Type.Object({
  name: MemberEntitySchema.properties.name,
  role: MemberEntitySchema.properties.role
})

export type CreateMemberDto = Static<typeof CreateMemberDtoSchema>
