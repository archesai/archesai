import type {
  Static,
  TBoolean,
  TNull,
  TNumber,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const API_TOKEN_ENTITY_KEY = 'api-tokens'

export const ApiTokenEntitySchema: TObject<{
  createdAt: TString
  enabled: TBoolean
  // use unions instaed of optional
  expiresAt: TUnion<[TString, TNull]>
  id: TString
  key: TString
  lastRefill: TUnion<[TString, TNull]>
  lastRequest: TUnion<[TString, TNull]>
  metadata: TUnion<[TString, TNull]>
  name: TUnion<[TString, TNull]>
  permissions: TUnion<[TString, TNull]>
  prefix: TUnion<[TString, TNull]>
  rateLimitEnabled: TBoolean
  rateLimitMax: TUnion<[TNumber, TNull]>
  rateLimitTimeWindow: TUnion<[TNumber, TNull]>
  refillAmount: TUnion<[TNumber, TNull]>
  refillInterval: TUnion<[TNumber, TNull]>
  remaining: TUnion<[TNumber, TNull]>
  requestCount: TNumber
  start: TUnion<[TString, TNull]>
  updatedAt: TString
  userId: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    enabled: Type.Boolean({
      description: 'Whether the API token is enabled or not'
    }),
    // use unions instaed of optional
    expiresAt: Type.Union([Type.String(), Type.Null()], {
      description: 'The date and time when the API token expires'
    }),
    key: Type.String({
      description: 'The API token key. This will only be shown once'
    }),
    lastRefill: Type.Union([Type.String(), Type.Null()], {
      description: 'The date and time when the API token was last refilled'
    }),
    lastRequest: Type.Union([Type.String(), Type.Null()], {
      description: 'The date and time when the API token was last used'
    }),
    metadata: Type.Union([Type.String(), Type.Null()], {
      description: 'The metadata for the API token, used for custom data'
    }),
    name: Type.Union([Type.String(), Type.Null()], {
      description: 'The name of the API token'
    }),
    permissions: Type.Union([Type.String(), Type.Null()], {
      description: 'The name of the API token'
    }),
    prefix: Type.Union([Type.String(), Type.Null()], {
      description: 'TThe prefix for the API token, used for routing requests'
    }),
    rateLimitEnabled: Type.Boolean({
      description: 'Whether the API token has rate limiting enabled'
    }),
    rateLimitMax: Type.Union([Type.Number(), Type.Null()], {
      description: 'The maximum number of requests allowed per time window'
    }),
    rateLimitTimeWindow: Type.Union([Type.Number(), Type.Null()], {
      description: 'The time window in seconds for the rate limit'
    }),
    refillAmount: Type.Union([Type.Number(), Type.Null()], {
      description: 'The amount of requests to refill the token with'
    }),
    refillInterval: Type.Union([Type.Number(), Type.Null()], {
      description: 'The interval in seconds to refill the token'
    }),
    remaining: Type.Union([Type.Number(), Type.Null()], {
      description: 'The number of requests remaining for the token'
    }),
    requestCount: Type.Number({
      description: 'The number of requests made with the token'
    }),
    start: Type.Union([Type.String(), Type.Null()], {
      description: 'The number of requests remaining for the token'
    }),
    userId: Type.String({
      description: 'The id of the user the token belongs to'
    })
  },
  {
    $id: 'ApiTokenEntity',
    description: 'The API token entity',
    title: 'API Token Entity'
  }
)

export type ApiTokenEntity = Static<typeof ApiTokenEntitySchema>
