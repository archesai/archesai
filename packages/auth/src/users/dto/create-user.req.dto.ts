import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '@archesai/schemas'

export const CreateUserRequestSchema = Type.Object({
  email: UserEntitySchema.properties.email,
  image: UserEntitySchema.properties.image,
  name: UserEntitySchema.properties.name,
  orgname: UserEntitySchema.properties.orgname
})

export type CreateUserRequest = Static<typeof CreateUserRequestSchema>
