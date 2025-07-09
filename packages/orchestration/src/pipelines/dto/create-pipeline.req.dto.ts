import { Type } from '@sinclair/typebox'

import { PipelineEntitySchema } from '@archesai/schemas'

export const CreatePipelineRequestSchema = Type.Object({
  description: PipelineEntitySchema.properties.description,
  name: PipelineEntitySchema.properties.name,
  steps: PipelineEntitySchema.properties.steps
})
