import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ConflictResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
      detail: Type.String({
        examples: [
          'The request conflicts with the current state of the resource.'
        ]
      }),
      status: Type.String({
        examples: ['409']
      }),
      title: Type.String({
        examples: ['Conflict']
      })
    })
  },
  {
    $id: 'ConflictResponse',
    description: 'Conflict',
    title: 'Conflict'
  }
)
