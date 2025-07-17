import type {
  Static,
  TNull,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateRunDtoSchema } from '#orchestration/runs/dto/create-run.dto'

export const UpdateRunDtoSchema: TObject<{
  pipelineId: TOptional<TUnion<[TString, TNull]>>
}> = Type.Partial(CreateRunDtoSchema)

export type UpdateRunDto = Static<typeof UpdateRunDtoSchema>
