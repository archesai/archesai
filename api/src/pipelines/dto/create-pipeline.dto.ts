import {
  ApiProperty,
  IntersectionType,
  PartialType,
  PickType,
} from "@nestjs/swagger";
import { IsArray, ValidateNested } from "class-validator";

import { PipelineEntity } from "../entities/pipeline.entity";
import { PipelineStepEntity } from "../entities/pipeline-step.entity";

export class CreatePipelineStepDto extends IntersectionType(
  PickType(PipelineStepEntity, ["toolId"] as const),
  PartialType(PickType(PipelineStepEntity, ["dependsOnId"] as const))
) {}

export class CreatePipelineDto extends PickType(PipelineEntity, [
  "name",
  "description",
]) {
  @ApiProperty({
    description: "An array of pipeline tools to be added to the pipeline",
    type: [CreatePipelineStepDto],
  })
  @IsArray()
  @ValidateNested({ each: true })
  pipelineSteps: CreatePipelineStepDto[];
}
