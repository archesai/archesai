import type { TArray, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const NotFoundResponseSchema: TObject<{
  errors: TArray<
    TObject<{
      detail: TString
      status: TString
      title: TString
    }>
  >
}> = Type.Object(
  {
    errors: Type.Array(
      Type.Object({
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
    )
  },
  {
    $id: 'NotFoundResponse',
    description: 'Not Found',
    title: 'Not Found'
  }
)
