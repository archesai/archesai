import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinkSchema } from '#http/schemas/links.schema'
import { MetaSchema } from '#http/schemas/meta.schema'
import { RelationshipsSchema } from '#http/schemas/relationship.schema'
import { toTitleCaseNoSpaces } from '#utils/strings'

export const createResourceObjectSchema = (
  entitySchema: TObject,
  entityKey: string,
  relationshipsSchema?: TObject
) => {
  if (!entityKey) {
    throw new Error('Entity schema must have an $id property')
  }
  return Type.Object(
    {
      attributes: Type.Omit(entitySchema, ['id']),
      id: Type.Pick(entitySchema, ['id']).properties.id,
      lid: Type.Optional(Type.String()),
      links: Type.Optional(
        Type.Object(
          {
            describedby: Type.Optional(LegacyRef(LinkSchema)),
            self: Type.Optional(LegacyRef(LinkSchema))
          },
          { additionalProperties: LinkSchema }
        )
      ),
      meta: Type.Optional(MetaSchema),
      relationships:
        relationshipsSchema ?
          Type.Object(relationshipsSchema)
        : Type.Optional(RelationshipsSchema),
      type: Type.Literal(entityKey)
    },
    {
      $id: `${toTitleCaseNoSpaces(entityKey)}Resource`,
      description: `${toTitleCaseNoSpaces(entityKey)} resource`,
      title: `${toTitleCaseNoSpaces(entityKey)} Resource`
    }
  )
}
