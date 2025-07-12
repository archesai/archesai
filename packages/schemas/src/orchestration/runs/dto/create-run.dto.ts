import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { RunEntitySchema } from '#orchestration/runs/entities/run.entity'

export const CreateRunDtoSchema = Type.Object({
  pipelineId: RunEntitySchema.properties.pipelineId
})

export type CreateRunDto = Static<typeof CreateRunDtoSchema>
