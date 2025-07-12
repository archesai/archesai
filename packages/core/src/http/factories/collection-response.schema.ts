import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinksSchema } from '#http/schemas/links.schema'
// import { MetaSchema } from '#http/schemas/meta.schema'
import { ResourceObjectSchema } from '#http/schemas/resource-object.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createCollectionResponseSchema = (
  resourceObjectSchema: TObject,
  entityKey: string
) => {
  return Type.Object(
    {
      data: Type.Array(resourceObjectSchema),
      included: Type.Optional(Type.Array(LegacyRef(ResourceObjectSchema))),
      links: Type.Optional(LegacyRef(LinksSchema))
      // meta: Type.Optional(LegacyRef(MetaSchema))
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}CollectionResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} collection response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Collection Response`
    }
  )
}
