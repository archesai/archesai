import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateUserRequestSchema } from '#users/dto/create-user.req.dto'

export const UpdateUserRequestSchema = Type.Partial(CreateUserRequestSchema)

export type UpdateUserRequest = Static<typeof UpdateUserRequestSchema>
