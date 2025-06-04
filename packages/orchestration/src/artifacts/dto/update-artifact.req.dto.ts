import { Type } from '@sinclair/typebox'

import { CreateArtifactRequestSchema } from '#artifacts/dto/create-artifact.req.dto'

export const UpdateArtifactRequestSchema = Type.Partial(
  CreateArtifactRequestSchema
)
