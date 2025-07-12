import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateMemberDtoSchema } from '#members/dto/create-member.dto'

export const UpdateMemberDtoSchema = Type.Partial(CreateMemberDtoSchema)

export type UpdateMemberDto = Static<typeof UpdateMemberDtoSchema>
