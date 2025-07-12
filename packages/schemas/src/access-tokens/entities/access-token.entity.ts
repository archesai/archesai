import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const AccessTokenEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    accessToken: Type.String({
      description:
        'The authorization token that can be used to access Arches AI'
    }),
    refreshToken: Type.String({
      description:
        'The refresh token that can be used to get a new access token'
    })
  },
  {
    $id: 'AccessTokenEntity',
    description: 'The access token entity',
    title: 'Access Token Entity'
  }
)

export type AccessTokenEntity = Static<typeof AccessTokenEntitySchema>

export const ACCESS_TOKEN_ENTITY_KEY = 'access-tokens'
