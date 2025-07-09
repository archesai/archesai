import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'

export const ApiTokenRoleTypes = ['ADMIN', 'USER'] as const
export type ApiTokenRoleType = (typeof ApiTokenRoleTypes)[number]

export const ApiTokenEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    key: Type.String({
      description: 'The API token key. This will only be shown once'
    }),
    orgname: Type.String({
      description: 'The name of the organization the token belongs to'
    }),
    role: Type.Union(
      ApiTokenRoleTypes.map((role) => Type.Literal(role)), // Using literals instead of enums
      { description: 'The role of the API token' }
    )
  },
  {
    $id: 'ApiTokenEntity',
    description: 'The API token entity',
    title: 'API Token Entity'
  }
)

export class ApiTokenEntity
  extends BaseEntity
  implements Static<typeof ApiTokenEntitySchema>
{
  public key: string
  public orgname: string
  public role: ApiTokenRoleType
  public type: string = API_TOKEN_ENTITY_KEY

  constructor(props: ApiTokenEntity) {
    super(props)
    this.key = props.key
    this.orgname = props.orgname
    this.role = props.role
  }
}

export const API_TOKEN_ENTITY_KEY = 'api-tokens'
