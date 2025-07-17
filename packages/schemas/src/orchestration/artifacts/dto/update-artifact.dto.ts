import type {
  Static,
  TNull,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreateArtifactDtoSchema } from '#orchestration/artifacts/dto/create-artifact.dto'

export const UpdateArtifactDtoSchema: TObject<{
  name: TOptional<TString>
  text: TOptional<TUnion<[TNull, TString]>>
  url: TOptional<TUnion<[TNull, TString]>>
}> = Type.Partial(CreateArtifactDtoSchema)

export type UpdateArtifactDto = Static<typeof UpdateArtifactDtoSchema>
