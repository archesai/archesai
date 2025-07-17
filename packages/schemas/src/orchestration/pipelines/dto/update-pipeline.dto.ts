import type {
  Static,
  TArray,
  TNull,
  TObject,
  TOptional,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreatePipelineDtoSchema } from '#orchestration/pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema: TObject<{
  description: TOptional<TUnion<[TString, TNull]>>
  name: TOptional<TUnion<[TString, TNull]>>
  steps: TOptional<
    TArray<
      TObject<{
        createdAt: TString
        dependents: TArray<
          TObject<{
            pipelineStepId: TString
          }>
        >
        id: TString
        pipelineId: TString
        prerequisites: TArray<
          TObject<{
            pipelineStepId: TString
          }>
        >
        toolId: TString
        updatedAt: TString
      }>
    >
  >
}> = Type.Partial(CreatePipelineDtoSchema)

export type UpdatePipelineDto = Static<typeof UpdatePipelineDtoSchema>
