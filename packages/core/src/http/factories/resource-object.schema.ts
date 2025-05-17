import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LinksSchema } from '#http/schemas/links.schema'
import { MetaSchema } from '#http/schemas/meta.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createResourceObjectSchema = (
  EntitySchema: TObject,
  entityKey: string
) => {
  return Type.Object(
    {
      attributes: Type.Omit(EntitySchema, ['id', 'type']),
      id: Type.Pick(EntitySchema, ['id']).properties.id,
      links: Type.Optional(Type.Pick(LinksSchema, ['self'])),
      meta: Type.Optional(MetaSchema),
      relationships: Type.Optional(Type.Object({})),
      type: Type.Literal(entityKey)
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}Resource`,
      description: `${toTitleCaseNoSpaces(entityKey)} resource`,
      title: `${toTitleCaseNoSpaces(entityKey)} Resource`
    }
  )
}
