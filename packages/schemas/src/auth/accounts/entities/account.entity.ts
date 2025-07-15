import type { Static, TNull, TObject, TString, TUnion } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const AccountEntitySchema: TObject<{
  accessToken: TUnion<[TString, TNull]>
  accessTokenExpiresAt: TUnion<[TString, TNull]>
  accountId: TString
  createdAt: TString
  id: TString
  idToken: TUnion<[TString, TNull]>
  password: TUnion<[TString, TNull]>
  providerId: TString
  refreshToken: TUnion<[TString, TNull]>
  refreshTokenExpiresAt: TUnion<[TString, TNull]>
  scope: TUnion<[TString, TNull]>
  updatedAt: TString
  userId: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    accessToken: Type.Union([Type.String(), Type.Null()], {
      description: 'The access token'
    }),
    accessTokenExpiresAt: Type.Union([Type.String(), Type.Null()], {
      description: 'The expiration date'
    }),
    accountId: Type.String({
      description: 'The unique identifier for the account'
    }),
    idToken: Type.Union([Type.String(), Type.Null()], {
      description: 'The ID token'
    }),
    password: Type.Union([Type.String(), Type.Null()], {
      description: 'The hashed password for local authentication'
    }),
    providerId: Type.String({
      description: 'The provider ID associated with the auth provider'
    }),
    refreshToken: Type.Union([Type.String(), Type.Null()], {
      description: 'The refresh token'
    }),
    refreshTokenExpiresAt: Type.Union([Type.String(), Type.Null()], {
      description: 'The refresh token expiration date'
    }),
    scope: Type.Union([Type.String(), Type.Null()], {
      description: 'The scope of the access token'
    }),
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
