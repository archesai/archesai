import { IntersectionType, PickType } from '@nestjs/swagger'
import { IsArray, IsString } from 'class-validator'

import { PipelineStepEntity } from '../entities/pipeline-step.entity'
import { PipelineEntity } from '../entities/pipeline.entity'
import { Expose } from 'class-transformer'

export class CreatePipelineDto extends PickType(PipelineEntity, [
  'name',
  'description',
  'pipelineSteps'
]) {}

export class CreatePipelineStepDto extends IntersectionType(
  PickType(PipelineStepEntity, ['toolId', 'id', 'name'] as const)
) {
  /**
   * An array of steps that this step depends on
   * @example ['step-id', 'step-id-2']
   */
  @IsArray()
  @IsString({ each: true })
  @Expose()
  dependsOn: string[]
}
