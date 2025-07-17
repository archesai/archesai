import type {
  Static,
  TArray,
  TObject,
  TOptional,
  TString
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { PipelineEntitySchema } from '#orchestration/pipelines/entities/pipeline.entity'

export const CreatePipelineDtoSchema: TObject<{
  description: TOptional<TString>
  name: TOptional<TString>
  steps: TArray<
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
}> = Type.Object({
  description: PipelineEntitySchema.properties.description,
  name: PipelineEntitySchema.properties.name,
  steps: PipelineEntitySchema.properties.steps
})

export type CreatePipelineDto = Static<typeof CreatePipelineDtoSchema>
