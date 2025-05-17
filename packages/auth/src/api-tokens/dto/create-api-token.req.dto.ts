import { Type } from '@sinclair/typebox'

import { ApiTokenEntitySchema, BaseEntitySchema } from '@archesai/domain'

export const CreateApiTokenRequestSchema = Type.Object({
  name: BaseEntitySchema.properties.name,
  orgname: ApiTokenEntitySchema.properties.orgname,
  role: ApiTokenEntitySchema.properties.role
})
