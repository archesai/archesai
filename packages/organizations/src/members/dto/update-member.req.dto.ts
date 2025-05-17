import { Type } from '@sinclair/typebox'

import { CreateMemberRequestSchema } from '#members/dto/create-member.req.dto'

export const UpdateMemberRequestSchema = Type.Partial(CreateMemberRequestSchema)
