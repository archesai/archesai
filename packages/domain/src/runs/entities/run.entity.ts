import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { RunType } from '#enums/role'

import { BaseEntity, BaseEntitySchema } from '#base/entities/base.entity'
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
        description: 'The timestamp when the run completed',
        format: 'date-time'
      })
    ),
    error: Type.Optional(Type.String({ description: 'The error message' })),
    orgname: Type.String({ description: 'The organization name' }),
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
        description: 'The timestamp when the run started',
        format: 'date-time'
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

export class RunEntity
  extends BaseEntity
  implements Static<typeof RunEntitySchema>
{
  // inputs: BaseEntity[]
  // outputs: BaseEntity[]
  // pipeline: BaseEntity
  // tool: BaseEntity
  public completedAt?: string
  public error?: string
  public orgname: string
  public pipelineId: string
  public progress: number
  // relationships: RunRelationships
  public runType: RunType
  public startedAt?: string
  public status: string
  public toolId: string
  public type = RUN_ENTITY_KEY

  constructor(props: RunEntity) {
    super(props)
    // this.inputs = props.inputs
    // this.outputs = props.outputs
    // this.pipeline = props.pipeline
    // this.tool = props.tool
    if (props.completedAt) {
      this.completedAt = props.completedAt
    }
    if (props.error) {
      this.error = props.error
    }
    if (props.startedAt) {
      this.startedAt = props.startedAt
    }
    this.orgname = props.orgname
    this.pipelineId = props.pipelineId
    this.progress = props.progress
    // this.relationships = props.relationships
    this.runType = props.runType
    this.status = props.status
    this.toolId = props.toolId
  }
}

export const RUN_ENTITY_KEY = 'runs'
