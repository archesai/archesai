import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '#users/entities/user.entity'

export const CreateUserDtoSchema = Type.Object({
  email: UserEntitySchema.properties.email,
  image: UserEntitySchema.properties.image,
  name: UserEntitySchema.properties.name,
  orgname: UserEntitySchema.properties.orgname
})

export type CreateUserDto = Static<typeof CreateUserDtoSchema>
