import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const BadRequestResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
      detail: Type.String({
        examples: ['The request is invalid or malformed.']
      }),
      status: Type.String({
        examples: ['400']
      }),
      title: Type.String({
        examples: ['Bad Request']
      })
    })
  },
  {
    $id: 'BadRequestResponse',
    description: 'Bad Request',
    title: 'Bad Request'
  }
)
