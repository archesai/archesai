import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const NotFoundResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
      detail: Type.String({
        examples: ['The requested resource could not be found.']
      }),
      status: Type.String({
        examples: ['404']
      }),
      title: Type.String({
        examples: ['Not Found']
      })
    })
  },
  {
    $id: 'NotFoundResponse',
    description: 'Not Found',
    title: 'Not Found'
  }
)
