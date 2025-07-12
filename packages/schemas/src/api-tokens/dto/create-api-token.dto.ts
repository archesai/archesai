import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

// import { BaseEntitySchema } from '@archesai/schemas'

export const CreateApiTokenDtoSchema = Type.Object({
  // name: BaseEntitySchema.properties.name
  // orgname: ApiTokenEntitySchema.properties.orgname,
  // role: ApiTokenEntitySchema.properties.role
})

export type CreateApiTokenDto = Static<typeof CreateApiTokenDtoSchema>
