import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinkSchema } from '#http/schemas/links.schema'
import { MetaSchema } from '#http/schemas/meta.schema'
import { RelationshipsSchema } from '#http/schemas/relationship.schema'
import { ResourceIdentifierSchema } from '#http/schemas/resource-identifier.schema'

export const ResourceObjectSchema = Type.Object(
  {
    ...ResourceIdentifierSchema.properties,
    attributes: Type.Optional(Type.Record(Type.String(), Type.Unknown())),
    links: Type.Optional(
      Type.Object(
        {
          describedby: Type.Optional(LegacyRef(LinkSchema)),
          self: Type.Optional(LegacyRef(LinkSchema))
        },
        { additionalProperties: LinkSchema }
      )
    ),
    meta: Type.Optional(LegacyRef(MetaSchema)),
    relationships: Type.Optional(LegacyRef(RelationshipsSchema))
  },
  {
    $id: 'ResourceObject',
    description: 'Resource object',
    title: 'Resource Object'
  }
)
