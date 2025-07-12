import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ErrorObjectSchema = Type.Object(
  {
    detail: Type.String({
      examples: ['The requested resource does not exist.']
    }),
    status: Type.String({
      examples: ['404']
    }),
    title: Type.String({
      examples: ['Not Found']
    })
  },
  {
    $id: 'ErrorObject',
    description: 'A list of errors that occurred during the request',
    title: 'Error Object'
  }
)

export type ErrorObject = Static<typeof ErrorObjectSchema>
