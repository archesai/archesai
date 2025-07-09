import { Type } from '@sinclair/typebox'

import { ArtifactEntitySchema } from '@archesai/schemas'

export const CreateArtifactRequestSchema = Type.Object({
  name: ArtifactEntitySchema.properties.name,
  text: ArtifactEntitySchema.properties.text,
  url: ArtifactEntitySchema.properties.url
})
