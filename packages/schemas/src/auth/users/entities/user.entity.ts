import type {
  Static,
  TBoolean,
  TNull,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const UserEntitySchema: TObject<{
  createdAt: TString
  deactivated: TBoolean
  email: TString
  emailVerified: TBoolean
  id: TString
  image: TUnion<[TNull, TString]>
  name: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    deactivated: Type.Boolean({
      default: false,
      description: 'Whether or not the user is deactivated'
    }),
    email: Type.String({
      description: "The user's e-mail"
      // format: 'email'
    }),
    emailVerified: Type.Boolean({
      description: "Whether or not the user's e-mail has been verified"
    }),
    image: Type.Union([Type.Null(), Type.String()], {
      description: "The user's avatar image URL"
    }),
    name: Type.String({
      description: "The user's name",
      minLength: 1
    })
  },
  {
    $id: 'UserEntity',
    description: 'The user entity',
    title: 'User Entity'
  }
)

export type UserEntity = Static<typeof UserEntitySchema>

export const USER_ENTITY_KEY = 'users'
