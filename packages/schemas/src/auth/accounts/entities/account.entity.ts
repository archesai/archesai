import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const AccountEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    accessToken: Type.Optional(
      Type.String({ description: 'The access token' })
    ),
    accessTokenExpiresAt: Type.Optional(
      Type.String({ description: 'The expiration date' })
    ),
    accountId: Type.String({
      description: 'The unique identifier for the account'
    }),
    idToken: Type.Optional(Type.String({ description: 'The ID token' })),
    password: Type.Optional(
      Type.String({
        description: 'The hashed password for local authentication'
      })
    ),
    providerId: Type.String({
      description: 'The provider ID associated with the auth provider'
    }),
    refreshToken: Type.Optional(
      Type.String({ description: 'The refresh token' })
    ),
    refreshTokenExpiresAt: Type.Optional(
      Type.String({ description: 'The refresh token expiration date' })
    ),
    scope: Type.Optional(
      Type.String({ description: 'The scope of the access token' })
    ),
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
