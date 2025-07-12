import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const UserEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    deactivated: Type.Boolean({
      default: false,
      description: 'Whether or not the user is deactivated'
    }),
    email: Type.String({
      description: "The user's e-mail",
      format: 'email'
    }),
    emailVerified: Type.Boolean({
      description: "Whether or not the user's e-mail has been verified",
      format: 'date-time'
    }),
    image: Type.Optional(
      Type.String({
        description: "The user's avatar image URL"
      })
    ),
    name: Type.String({
      description: "The user's name",
      minLength: 1
    }),
    orgname: Type.String({
      description: "The user's organization name"
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
