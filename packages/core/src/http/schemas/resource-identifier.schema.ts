import { Type } from '@sinclair/typebox'

export const ResourceIdentifierSchema = Type.Object(
  {
    id: Type.String(),
    type: Type.String()
  },
  {
    $id: 'ResourceIdentifier',
    description: 'Resource Identifier',
    title: 'Resource Identifier'
  }
)
