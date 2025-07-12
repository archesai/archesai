import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
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

export type AccountEntity = Static<typeof AccountEntitySchema>

export const ACCOUNT_ENTITY_KEY = 'accounts'
