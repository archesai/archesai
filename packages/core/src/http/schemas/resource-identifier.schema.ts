import { Type } from '@sinclair/typebox'

export const ResourceIdentifierSchema = Type.Object(
  {
    id: Type.String({
      description: 'Unique identifier for the resource'
    }),
    type: Type.String({
      description: 'Type of the resource'
    })
  },
  {
    $id: 'ResourceIdentifier',
    description: 'Resource Identifier',
    title: 'Resource Identifier'
  }
)
