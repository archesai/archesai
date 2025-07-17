import type { TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const ForbiddenResponseSchema: TObject<{
  error: TObject<{
    detail: TString
    status: TString
    title: TString
  }>
}> = Type.Object(
  {
    error: Type.Object({
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
  },
  {
    $id: 'ForbiddenResponse',
    description: 'Forbidden',
    title: 'Forbidden'
  }
)
