import { SubItemEntity } from '@/src/common/entities/base-sub-item.entity'
import { BaseEntity } from '@/src/common/entities/base.entity'
import { ToolEntity } from '@/src/tools/entities/tool.entity'
import {
  PipelineStep as _PrismaPipelineStep,
  Tool as _PrismaTool
} from '@prisma/client'
import { Expose } from 'class-transformer'
import { IsString, ValidateNested } from 'class-validator'

type PipelineStepModel = _PrismaPipelineStep & {
  tool: _PrismaTool
}

export class PipelineStepEntity
  extends BaseEntity
  implements PipelineStepModel
{
  /**
   * The order of the step in the pipeline
   */
  @ValidateNested({ each: true })
  @Expose()
  dependents: SubItemEntity[]

  /**
   * These are the steps that this step depends on.
   */
  @ValidateNested({ each: true })
  @Expose()
  dependsOn: SubItemEntity[]

  /**
   * The name of the step in the pipeline. It must be unique within the pipeline.
   */
  @IsString()
  @Expose()
  name: string

  /**
   * The ID of the pipelin that this step belongs to
   * @example 'pipeline-id'
   */
  @IsString()
  @Expose()
  pipelineId: string

  /**
   * The name of the tool that this step uses.
   */
  @ValidateNested()
  @Expose()
  tool: ToolEntity

  /**
   * This is the ID of the tool that this step uses.
   * @example 'tool-id'
   */
  @IsString()
  @Expose()
  toolId: string

  constructor(pipelineStep: PipelineStepModel) {
    super()
    Object.assign(this, pipelineStep)
  }
}
