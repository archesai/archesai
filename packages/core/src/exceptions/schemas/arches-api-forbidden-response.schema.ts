import { Type } from '@sinclair/typebox'

export const ArchesApiForbiddenResponseSchema = Type.Object(
  {
    errors: Type.Array(
      Type.Object({
        detail: Type.String({
          examples: ['You do not have permission to access this resource.']
        }),
        status: Type.String({
          examples: ['403']
        }),
        title: Type.String({
          examples: ['Forbidden']
        })
      })
    )
  },
  {
    $id: 'ForbiddenResponse',
    description: 'Forbidden',
    title: 'Forbidden'
  }
)
