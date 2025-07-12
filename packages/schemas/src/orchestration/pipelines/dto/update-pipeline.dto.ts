import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreatePipelineDtoSchema } from '#orchestration/pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema = Type.Partial(CreatePipelineDtoSchema)

export type UpdatePipelineDto = Static<typeof UpdatePipelineDtoSchema>
