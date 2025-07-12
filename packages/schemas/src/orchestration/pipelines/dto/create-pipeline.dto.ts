import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { PipelineEntitySchema } from '#orchestration/pipelines/entities/pipeline.entity'

export const CreatePipelineDtoSchema = Type.Object({
  description: PipelineEntitySchema.properties.description,
  name: PipelineEntitySchema.properties.name,
  steps: PipelineEntitySchema.properties.steps
})

export type CreatePipelineDto = Static<typeof CreatePipelineDtoSchema>
