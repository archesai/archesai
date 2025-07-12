import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinksSchema } from '#http/schemas/links.schema'
// import { MetaSchema } from '#http/schemas/meta.schema'
import { ResourceObjectSchema } from '#http/schemas/resource-object.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createIndividualResponseSchema = (
  resourceObjectSchema: TObject,
  entityKey: string
) => {
  return Type.Object(
    {
      data: resourceObjectSchema,
      included: Type.Optional(Type.Array(LegacyRef(ResourceObjectSchema))),
      links: Type.Optional(LegacyRef(LinksSchema))
      // meta: Type.Optional(LegacyRef(MetaSchema))
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}IndividualResponse`,
      description: `${toTitleCaseNoSpaces(entityKey)} Individual response`,
      title: `${toTitleCaseNoSpaces(entityKey)} Individual Response`
    }
  )
}
