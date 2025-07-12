import { Type } from '@sinclair/typebox'

import { ResourceIdentifierSchema } from '#http/schemas/resource-identifier.schema'

const RelationshipSchema = Type.Object(
  {
    data: ResourceIdentifierSchema
  },
  {
    $id: 'Relationship',
    description: 'Relationship object',
    title: 'Relationship'
  }
)

export const RelationshipsSchema = Type.Record(
  Type.String(),
  RelationshipSchema,
  {
    $id: 'Relationships',
    description: 'Relationships object',
    title: 'Relationships'
  }
)
