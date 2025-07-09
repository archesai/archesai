import { Type } from '@sinclair/typebox'

import { RunEntitySchema } from '@archesai/schemas'

export const CreateRunRequestSchema = Type.Object({
  pipelineId: RunEntitySchema.properties.pipelineId
})
