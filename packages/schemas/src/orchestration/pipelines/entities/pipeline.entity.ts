import type {
  Static,
  TArray,
  TObject,
  TOptional,
  TString
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const PipelineStepEntitySchema: TObject<{
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
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    dependents: Type.Array(
      Type.Object({
        pipelineStepId: Type.String()
      })
    ),
    pipelineId: Type.String(),
    prerequisites: Type.Array(
      Type.Object({
        pipelineStepId: Type.String()
      })
    ),
    toolId: Type.String()
  },
  {
    $id: 'PipelineStepEntity',
    description: 'The pipeline step entity',
    title: 'Pipeline Step Entity'
  }
)

export const PipelineEntitySchema: TObject<{
  createdAt: TString
  description: TOptional<TString>
  id: TString
  name: TOptional<TString>
  organizationId: TString
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
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    description: Type.Optional(
      Type.String({ description: 'The pipeline description' })
    ),
    name: Type.Optional(Type.String({ description: 'The pipeline name' })),
    organizationId: Type.String({ description: 'The organization id' }),
    steps: Type.Array(PipelineStepEntitySchema, {
      description: 'The steps in the pipeline'
    })
  },
  {
    $id: 'PipelineEntity',
    description: 'The pipeline entity',
    title: 'Pipeline Entity'
  }
)

export type PipelineEntity = Static<typeof PipelineEntitySchema>

export type PipelineStepEntity = Static<typeof PipelineStepEntitySchema>

export const PIPELINE_ENTITY_KEY = 'pipelines'

export const PIPELINE_STEP_ENTITY_KEY = 'pipeline-steps'
