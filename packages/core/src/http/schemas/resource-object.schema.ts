import { Type } from '@sinclair/typebox'

import { LegacyRef } from '@archesai/schemas'

import { LinkSchema } from '#http/schemas/links.schema'
import { RelationshipsSchema } from '#http/schemas/relationship.schema'
import { ResourceIdentifierSchema } from '#http/schemas/resource-identifier.schema'

export const ResourceObjectSchema = Type.Object(
  {
    ...ResourceIdentifierSchema.properties,
    attributes: Type.Optional(Type.Record(Type.String(), Type.Unknown())),
    links: Type.Optional(
      Type.Object({
        self: Type.Optional(LegacyRef(LinkSchema))
      })
    ),
    relationships: Type.Optional(LegacyRef(RelationshipsSchema))
  },
  {
    $id: 'ResourceObject',
    description: 'Resource object',
    title: 'Resource Object'
  }
)
