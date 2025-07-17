import type {
  Static,
  TLiteral,
  TNull,
  TNumber,
  TObject,
  TString,
  TUnion
} from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { StatusTypes } from '#enums/role'

// export const runRelationshipsSchema = object({
//   inputs: z.array(baseEntitySchema).describe('The inputs to the run'),
//   outputs: z.array(baseEntitySchema).describe('The outputs of the run'),
//   pipeline: baseEntitySchema
//     .nullable()
//     .describe('The pipeline associated with the run'),
//   tool: baseEntitySchema.nullable().describe('The tool associated with the run')
// })

export const RunEntitySchema: TObject<{
  completedAt: TUnion<[TString, TNull]>
  createdAt: TString
  error: TUnion<[TString, TNull]>
  id: TString
  organizationId: TString
  pipelineId: TUnion<[TString, TNull]>
  progress: TNumber
  startedAt: TUnion<[TString, TNull]>
  status: TUnion<TLiteral<'COMPLETED' | 'FAILED' | 'PROCESSING' | 'QUEUED'>[]>
  toolId: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    completedAt: Type.Union([Type.String(), Type.Null()], {
      description: 'The timestamp when the run completed'
    }),
    error: Type.Union([Type.String(), Type.Null()], {
      description: 'The error message'
    }),
    organizationId: Type.String({ description: 'The organization name' }),
    pipelineId: Type.Union([Type.String(), Type.Null()], {
      description: 'The pipeline ID associated with the run'
    }),
    progress: Type.Number({ description: 'The percent progress of the run' }),
    startedAt: Type.Union([Type.String(), Type.Null()], {
      description: 'The timestamp when the run started'
    }),
    status: Type.Union(StatusTypes.map((status) => Type.Literal(status))),
    toolId: Type.String({
      description: 'The tool ID associated with the run'
    })
  },
  {
    $id: 'RunEntity',
    description: 'The run entity',
    title: 'Run Entity'
  }
)

export type RunEntity = Static<typeof RunEntitySchema>

export const RUN_ENTITY_KEY = 'runs'
