import { Type } from '@sinclair/typebox'

import { ArtifactEntitySchema } from '@archesai/domain'

export const CreateArtifactRequestSchema = Type.Object({
  name: ArtifactEntitySchema.properties.name,
  text: ArtifactEntitySchema.properties.text,
  url: ArtifactEntitySchema.properties.url
})
