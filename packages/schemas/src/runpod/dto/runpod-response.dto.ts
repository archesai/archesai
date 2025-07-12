import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const RunpodResponseDtoSchema = Type.Object({
  id: Type.String(),
  output: Type.String(),
  status: Type.Union([
    Type.Literal('IN_PROGRESS'),
    Type.Literal('COMPLETED'),
    Type.Literal('FAILED')
  ])
})

export type RunpodResponseDto = Static<typeof RunpodResponseDtoSchema>
