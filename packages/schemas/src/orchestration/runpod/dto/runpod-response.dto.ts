import type {
  Static,
  TLiteral,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const RunpodResponseDtoSchema: TObject<{
  id: TString
  output: TString
  status: TUnion<
    [TLiteral<'IN_PROGRESS'>, TLiteral<'COMPLETED'>, TLiteral<'FAILED'>]
  >
}> = Type.Object({
  id: Type.String(),
  output: Type.String(),
  status: Type.Union([
    Type.Literal('IN_PROGRESS'),
    Type.Literal('COMPLETED'),
    Type.Literal('FAILED')
  ])
})

export type RunpodResponseDto = Static<typeof RunpodResponseDtoSchema>
