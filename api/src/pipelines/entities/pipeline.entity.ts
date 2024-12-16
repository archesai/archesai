import { BaseEntity } from '@/src/common/entities/base.entity'
import {
  Pipeline as _PrismaPipeline,
  PipelineStep as _PrismaPipelineStep,
  Tool as _PrismaTool
} from '@prisma/client'
import { IsOptional, IsString, ValidateNested } from 'class-validator'

import { PipelineStepEntity } from './pipeline-step.entity'
import { Expose } from 'class-transformer'

export type PipelineWithPipelineStepsModel = _PrismaPipeline & {
  pipelineSteps: (_PrismaPipelineStep & {
    dependents: { id: string; name: string }[]
    dependsOn: { id: string; name: string }[]
    tool: _PrismaTool
  })[]
}

export class PipelineEntity
  extends BaseEntity
  implements PipelineWithPipelineStepsModel
{
  /**
   * The description of the pipeline
   * @example 'This pipeline does something'
   */
  @IsOptional()
  @IsString()
  @Expose()
  description: null | string

  /**
   * The name of the pipeline
   * @example 'my-pipeline'
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The name of the organization that this pipeline belongs to
   * @example 'my-org'
   */
  @IsString()
  @Expose()
  orgname: string

  /**
   * The steps in the pipeline
   */
  @ValidateNested({ each: true })
  @Expose()
  pipelineSteps: PipelineStepEntity[]

  constructor(pipeline: PipelineWithPipelineStepsModel) {
    super()
    Object.assign(this, pipeline)
    this.pipelineSteps = pipeline.pipelineSteps.map(
      (pipelineStep) => new PipelineStepEntity(pipelineStep)
    )
  }
}
