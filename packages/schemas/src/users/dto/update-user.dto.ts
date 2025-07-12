import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateUserDtoSchema } from '#users/dto/create-user.dto'

export const UpdateUserDtoSchema = Type.Partial(CreateUserDtoSchema)

export type UpdateUserDto = Static<typeof UpdateUserDtoSchema>
