import { Type } from '@sinclair/typebox'

// import { MetaSchema } from '#http/schemas/meta.schema'
import { ResourceIdentifierSchema } from '#http/schemas/resource-identifier.schema'

const RelationshipLinksSchema = Type.Optional(
  Type.Partial(
    Type.Object({
      related: Type.Optional(
        Type.String({
          description: 'The related page for the resource item',
          examples: ['https://api.example.com/v1/items/1/author']
        })
      ),
      self: Type.Optional(
        Type.String({
          description: 'The current page for the resource item',
          examples: ['https://api.example.com/v1/items/1/relationships/author']
        })
      )
    })
  )
)

export const RelationshipSchema = Type.Object({
  data: Type.Optional(
    Type.Union([
      ResourceIdentifierSchema,
      Type.Array(ResourceIdentifierSchema),
      Type.Null()
    ])
  ),
  links: RelationshipLinksSchema
  // meta: MetaSchema
})
