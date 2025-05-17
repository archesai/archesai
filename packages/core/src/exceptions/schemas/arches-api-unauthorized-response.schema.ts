import type { TSchema } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ArchesApiUnauthorizedResponseSchema = Type.Object(
  {
    errors: Type.Array(
      Type.Object({
        detail: Type.String({
          examples: ['You are not authrozied to reach this endpoint.']
        }),
        status: Type.String({
          examples: ['401']
        }),
        title: Type.String({
          examples: ['Unauthorized']
        })
      })
    )
  },
  {
    $id: 'UnauthorizedResponse',
    description: 'Unauthorized',
    title: 'Unauthorized'
  }
) satisfies TSchema as TSchema
