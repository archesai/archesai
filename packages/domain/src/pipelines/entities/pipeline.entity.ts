// filepath: /home/jonathan/Projects/archesai/packages/domain/src/pipelines/entities/pipeline.entity.ts
import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { ToolEntity } from '#tools/entities/tool.entity'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
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

export class PipelineEntity
  extends BaseEntity
  implements Static<typeof PipelineEntitySchema>
{
  public description: string
  public orgname: string
  public steps: PipelineStepEntity[]
  public type = PIPELINE_ENTITY_KEY

  constructor(props: PipelineEntity) {
    super(props)
    this.description = props.description
    this.orgname = props.orgname
    this.steps = props.steps
  }
}

export class PipelineStepEntity
  extends BaseEntity
  implements Static<typeof PipelineStepEntitySchema>
{
  public dependents: { pipelineStepId: string }[]
  public pipelineId: string
  public prerequisites: { pipelineStepId: string }[]
  public tool: ToolEntity
  public toolId: string
  public type = PIPELINE_STEP_ENTITY_KEY

  constructor(props: PipelineStepEntity) {
    super(props)
    this.dependents = props.dependents
    this.pipelineId = props.pipelineId
    this.prerequisites = props.prerequisites
    this.tool = props.tool
    this.toolId = props.toolId
  }
}

export const PIPELINE_ENTITY_KEY = 'pipelines'

export const PIPELINE_STEP_ENTITY_KEY = 'pipeline-steps'
