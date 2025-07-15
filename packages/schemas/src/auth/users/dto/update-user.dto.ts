import type {
  Static,
  TNull,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateUserDtoSchema } from '#auth/users/dto/create-user.dto'

export const UpdateUserDtoSchema: TObject<{
  email: TOptional<TString>
  image: TOptional<TUnion<[TNull, TString]>>
}> = Type.Partial(CreateUserDtoSchema)

export type UpdateUserDto = Static<typeof UpdateUserDtoSchema>
