import type { Static, TNull, TObject, TString, TUnion } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { UserEntitySchema } from '#auth/users/entities/user.entity'

export const CreateUserDtoSchema: TObject<{
  email: TString
  image: TUnion<[TNull, TString]>
}> = Type.Object({
  email: UserEntitySchema.properties.email,
  image: UserEntitySchema.properties.image
})

export type CreateUserDto = Static<typeof CreateUserDtoSchema>
