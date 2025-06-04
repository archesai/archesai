import { Type } from '@sinclair/typebox'

import { RunEntitySchema } from '@archesai/domain'

export const CreateRunRequestSchema = Type.Object({
  pipelineId: RunEntitySchema.properties.pipelineId
})
