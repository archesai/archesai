import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

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

export class UserEntity
  extends BaseEntity
  implements Static<typeof UserEntitySchema>
{
  public deactivated: boolean
  public email: string
  public emailVerified?: string
  public image?: string
  public orgname: string
  public type = USER_ENTITY_KEY

  constructor(props: UserEntity) {
    super(props)
    this.deactivated = props.deactivated
    this.email = props.email
    if (props.emailVerified) {
      this.emailVerified = props.emailVerified
    }
    if (props.image) {
      this.image = props.image
    }

    this.orgname = props.orgname
  }
}

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
