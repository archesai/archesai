import { Type } from '@sinclair/typebox'

// import { MetaSchema } from '#http/schemas/meta.schema'

export const ResourceIdentifierSchema = Type.Object(
  {
    id: Type.String({
      description: 'Resource unique identifier',
      examples: ['9']
    }),
    // meta: MetaSchema,
    type: Type.String({
      description: 'Resource type identifier',
      examples: ['people']
    })
  },
  {
    description: 'A JSON:API-compliant resource identifier object',
    title: 'Resource Identifier'
  }
)
