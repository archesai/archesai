import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const InternalServerErrorResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
      detail: Type.String({
        examples: ['An unexpected error occurred on the server.']
      }),
      status: Type.String({
        examples: ['500']
      }),
      title: Type.String({
        examples: ['Internal Server Error']
      })
    })
  },
  {
    $id: 'InternalServerErrorResponse',
    description: 'Internal Server Error',
    title: 'Internal Server Error'
  }
)
