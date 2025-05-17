import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { AuthType, ProviderType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
import { AuthTypes, ProviderTypes } from '#enums/role'

export const AccountEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    access_token: Type.Optional(
      Type.String({ description: 'The access token' })
    ),
    authType: Type.Union(
      AuthTypes.map((authType) => Type.Literal(authType)),
      {
        description: 'The type of auth provider'
      }
    ),
    expires_at: Type.Optional(
      Type.String({ description: 'The expiration date' })
    ),
    hashed_password: Type.Optional(
      Type.String({
        description: 'The hashed password for local authentication'
      })
    ),
    id_token: Type.Optional(Type.String({ description: 'The ID token' })),
    provider: Type.Union(
      ProviderTypes.map((provider) => Type.Literal(provider)),
      {
        description: 'The auth provider name'
      }
    ),
    providerAccountId: Type.String({
      description: 'The provider ID associated with the auth provider'
    }),
    refresh_token: Type.Optional(
      Type.String({ description: 'The refresh token' })
    ),
    scope: Type.Optional(
      Type.String({ description: 'The scope of the access token' })
    ),
    session_state: Type.Optional(
      Type.String({ description: 'The session state' })
    ),
    token_type: Type.Optional(Type.String({ description: 'The token type' })),
    userId: Type.String({
      description: 'The user ID associated with the auth provider'
    })
  },
  {
    $id: 'AccountEntity',
    description: 'The account entity',
    title: 'Account Entity'
  }
)

export class AccountEntity
  extends BaseEntity
  implements Static<typeof AccountEntitySchema>
{
  public access_token?: string
  public authType: AuthType
  public expires_at?: string
  public hashed_password?: string
  public id_token?: string
  public provider: ProviderType
  public providerAccountId: string
  public refresh_token?: string
  public scope?: string
  public session_state?: string
  public token_type?: string
  public type = ACCOUNT_ENTITY_KEY
  public userId: string

  constructor(props: AccountEntity) {
    super(props)
    this.provider = props.provider
    this.providerAccountId = props.providerAccountId
    this.userId = props.userId
    this.authType = props.authType
    if (props.expires_at) {
      this.expires_at = props.expires_at
    }
    if (props.id_token) {
      this.id_token = props.id_token
    }
    if (props.session_state) {
      this.session_state = props.session_state
    }
    if (props.token_type) {
      this.token_type = props.token_type
    }
    if (props.hashed_password) {
      this.hashed_password = props.hashed_password
    }
    if (props.access_token) {
      this.access_token = props.access_token
    }
    if (props.scope) {
      this.scope = props.scope
    }
    if (props.token_type) {
      this.token_type = props.token_type
    }
  }
}

export const ACCOUNT_ENTITY_KEY = 'accounts'
