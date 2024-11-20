import { IntersectionType, PickType } from "@nestjs/swagger";
import { IsArray, ValidateNested } from "class-validator";

import { PipelineEntity } from "../entities/pipeline.entity";
import { PipelineStepEntity } from "../entities/pipeline-step.entity";

export class CreatePipelineStepDto extends IntersectionType(
  PickType(PipelineStepEntity, ["toolId", "id"] as const)
) {
  /**
   * An array of steps that this step depends on
   * @example ['step-id', 'step-id-2']
   */
  @IsArray()
  dependsOn: string[];
}

export class CreatePipelineDto extends PickType(PipelineEntity, [
  "name",
  "description",
]) {
  /**
   * An array of pipeline tools to be added to the pipeline
   */
  @IsArray()
  @ValidateNested({ each: true })
  pipelineSteps: CreatePipelineStepDto[];
}
