import { IntersectionType, PickType } from '@nestjs/swagger'
import { IsArray, IsString, ValidateNested } from 'class-validator'

import { PipelineStepEntity } from '../entities/pipeline-step.entity'
import { PipelineEntity } from '../entities/pipeline.entity'
import { Expose } from 'class-transformer'

export class CreatePipelineDto extends PickType(PipelineEntity, [
  'name',
  'description'
]) {
  /**
   * An array of pipeline tools to be added to the pipeline
   */
  @IsArray()
  @ValidateNested({ each: true })
  @Expose()
  pipelineSteps: CreatePipelineStepDto[]
}

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
