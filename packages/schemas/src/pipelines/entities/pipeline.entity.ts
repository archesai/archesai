import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { ToolEntitySchema } from '#tools/entities/tool.entity'

export const PipelineStepEntitySchema = Type.Object(
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
    tool: ToolEntitySchema,
    toolId: Type.String()
  },
  {
    $id: 'PipelineStepEntity',
    description: 'The pipeline step entity',
    title: 'Pipeline Step Entity'
  }
)

export const PipelineEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    description: Type.String({ description: 'The pipeline description' }),
    orgname: Type.String({ description: 'The organization name' }),
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
