import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { RunTypes } from '#enums/role'

// export const runRelationshipsSchema = object({
//   inputs: z.array(baseEntitySchema).describe('The inputs to the run'),
//   outputs: z.array(baseEntitySchema).describe('The outputs of the run'),
//   pipeline: baseEntitySchema
//     .nullable()
//     .describe('The pipeline associated with the run'),
//   tool: baseEntitySchema.nullable().describe('The tool associated with the run')
// })

export const RunEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    completedAt: Type.Optional(
      Type.String({
        description: 'The timestamp when the run completed'
      })
    ),
    error: Type.Optional(Type.String({ description: 'The error message' })),
    organizationId: Type.String({ description: 'The organization name' }),
    pipelineId: Type.String({
      description: 'The pipeline ID associated with the run'
    }),
    progress: Type.Number({ description: 'The percent progress of the run' }),
    runType: Type.Union(
      RunTypes.map((type) => Type.Literal(type)),
      { description: 'The type of run' }
    ),
    startedAt: Type.Optional(
      Type.String({
        description: 'The timestamp when the run started'
      })
    ),
    status: Type.String({ description: 'The status of the run' }),
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
