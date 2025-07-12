import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreatePipelineDtoSchema } from '#pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema = Type.Partial(CreatePipelineDtoSchema)

export type UpdatePipelineDto = Static<typeof UpdatePipelineDtoSchema>
