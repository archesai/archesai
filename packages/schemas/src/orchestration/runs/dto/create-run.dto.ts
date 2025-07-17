import type { Static, TNull, TObject, TString, TUnion } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { RunEntitySchema } from '#orchestration/runs/entities/run.entity'

export const CreateRunDtoSchema: TObject<{
  pipelineId: TUnion<[TString, TNull]>
}> = Type.Object({
  pipelineId: RunEntitySchema.properties.pipelineId
})

export type CreateRunDto = Static<typeof CreateRunDtoSchema>
