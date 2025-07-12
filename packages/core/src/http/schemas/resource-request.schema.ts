import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { MetaSchema } from '#http/schemas/meta.schema'
import { ResourceObjectSchema } from '#http/schemas/resource-object.schema'

// Create/Update Request
export const ResourceRequest = Type.Object(
  {
    data: LegacyRef(ResourceObjectSchema),
    meta: Type.Optional(LegacyRef(MetaSchema))
  },
  {
    $id: 'ResourceRequest',
    description: 'Resource Request',
    title: 'Resource Request'
  }
)
