import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ArtifactEntitySchema } from '#artifacts/entities/artifact.entity'

export const CreateArtifactDtoSchema = Type.Object({
  name: ArtifactEntitySchema.properties.name,
  text: ArtifactEntitySchema.properties.text,
  url: ArtifactEntitySchema.properties.url
})

export type CreateArtifactDto = Static<typeof CreateArtifactDtoSchema>
