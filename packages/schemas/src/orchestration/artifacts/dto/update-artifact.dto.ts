import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateArtifactDtoSchema } from '#orchestration/artifacts/dto/create-artifact.dto'

export const UpdateArtifactDtoSchema = Type.Partial(CreateArtifactDtoSchema)

export type UpdateArtifactDto = Static<typeof UpdateArtifactDtoSchema>
