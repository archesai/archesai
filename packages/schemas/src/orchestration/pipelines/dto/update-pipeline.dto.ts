import type {
  Static,
  TArray,
  TObject,
  TOptional,
  TString
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { CreatePipelineDtoSchema } from '#orchestration/pipelines/dto/create-pipeline.dto'

export const UpdatePipelineDtoSchema: TObject<{
  description: TOptional<TString>
  name: TOptional<TString>
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
        name: TString
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
