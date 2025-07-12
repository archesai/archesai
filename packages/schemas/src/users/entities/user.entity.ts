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
    emailVerified: Type.Optional(
      Type.String({
        description: "Whether or not the user's e-mail has been verified",
        format: 'date-time'
      })
    ),
    image: Type.Optional(
      Type.String({
        description: "The user's avatar image URL"
      })
    ),
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

export class UserRelations {
  public accounts: { id: string; type: string }[]
  public memberships: { id: string; type: string }[]

  constructor(props: UserRelations) {
    this.accounts = props.accounts
    this.memberships = props.memberships
  }

  public static schema() {
    return Type.Object({
      accounts: Type.Array(
        Type.Object({
          id: Type.String(),
          type: Type.String()
        }),
        { description: 'The accounts associated with the user' }
      ),
      memberships: Type.Array(
        Type.Object({
          id: Type.String(),
          type: Type.String()
        }),
        { description: 'The memberships associated with the user' }
      )
    })
  }
}

export const USER_ENTITY_KEY = 'users'
