import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { ArtifactEntitySchema } from '#orchestration/artifacts/entities/artifact.entity'

export const CreateArtifactDtoSchema = Type.Object({
  name: Type.String({
    description: 'The name of the artifact'
  }),
  text: ArtifactEntitySchema.properties.text,
  url: ArtifactEntitySchema.properties.url
})

export type CreateArtifactDto = Static<typeof CreateArtifactDtoSchema>
