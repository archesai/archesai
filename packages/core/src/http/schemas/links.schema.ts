import { Type } from '@sinclair/typebox'

export const LinksSchema = Type.Object({
  first: Type.Optional(
    Type.String({
      description: `The first page for the paginated resource collection`,
      examples: [`https://api.example.com/v1/items?page=1`]
    })
  ),
  last: Type.Optional(
    Type.String({
      description: `The last page for the paginated resource collection`,
      examples: [`https://api.example.com/v1/items?page=3`]
    })
  ),
  next: Type.Optional(
    Type.String({
      description: `The next page for the paginated resource collection`,
      examples: [`https://api.example.com/v1/items?page=3`]
    })
  ),
  prev: Type.Optional(
    Type.String({
      description: `The previous page for the paginated resource collection`,
      examples: [`https://api.example.com/v1/items?page=1`]
    })
  ),
  self: Type.Optional(
    Type.String({
      description: `The current page for the paginated resource collection`,
      examples: [`https://api.example.com/v1/items?page=2`]
    })
  )
})
