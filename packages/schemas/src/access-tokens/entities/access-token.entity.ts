import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { BaseInsertion } from '#base/entities/base.entity'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

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

export class AccessTokenEntity
  extends BaseEntity
  implements Static<typeof AccessTokenEntitySchema>
{
  public accessToken: string
  public refreshToken: string
  public type = ACCESS_TOKEN_ENTITY_KEY

  constructor(props: BaseInsertion<AccessTokenEntity>) {
    super(props)
    this.accessToken = props.accessToken
    this.refreshToken = props.refreshToken
  }
}

export const ACCESS_TOKEN_ENTITY_KEY = 'access-tokens'
