import { Type } from '@sinclair/typebox'

import { CreatePipelineRequestSchema } from '#pipelines/dto/create-pipeline.req.dto'

export const UpdatePipelineRequestSchema = Type.Partial(
  CreatePipelineRequestSchema
)
