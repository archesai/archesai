import type { TObject } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinkSchema } from '#http/schemas/links.schema'
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
      links: Type.Optional(
        Type.Object({
          self: Type.Optional(LegacyRef(LinkSchema))
        })
      ),
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
