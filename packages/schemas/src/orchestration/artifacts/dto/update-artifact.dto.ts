import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateArtifactDtoSchema } from '#orchestration/artifacts/dto/create-artifact.dto'

export const UpdateArtifactDtoSchema: TObject<{
  name: TOptional<TString>
  text: TOptional<TString>
  url: TOptional<TString>
}> = Type.Partial(CreateArtifactDtoSchema)

export type UpdateArtifactDto = Static<typeof UpdateArtifactDtoSchema>
