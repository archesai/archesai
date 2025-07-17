import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ValidationErrorResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
      detail: Type.String({
        examples: ['Validation failed for one or more fields.']
      }),
      details: Type.Array(
        Type.Object({
          field: Type.String({
            examples: ['username', 'email']
          }),
          message: Type.String({
            examples: ['Username is required.', 'Email format is invalid.']
          }),
          value: Type.Optional(
            Type.String({
              examples: ['john_doe', 'invalid-email']
            })
          )
        })
      ),
      status: Type.String({
        examples: ['422']
      }),
      title: Type.String({
        examples: ['Validation Error']
      })
    })
  },
  {
    $id: 'ValidationErrorResponse',
    description: 'Validation Error',
    title: 'Validation Error'
  }
)
