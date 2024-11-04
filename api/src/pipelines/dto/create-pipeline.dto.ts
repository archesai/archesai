import {
  ApiProperty,
  IntersectionType,
  PartialType,
  PickType,
} from "@nestjs/swagger";
import { IsArray, ValidateNested } from "class-validator";

import { PipelineEntity } from "../entities/pipeline.entity";
import { PipelineToolEntity } from "../entities/pipeline-tool.entity";

export class CreatePipelineToolDto extends IntersectionType(
  PickType(PipelineToolEntity, ["toolId"] as const),
  PartialType(PickType(PipelineToolEntity, ["dependsOnId"] as const))
) {}

export class CreatePipelineDto extends PickType(PipelineEntity, [
  "name",
  "description",
]) {
  @ApiProperty({
    description: "An array of pipeline tools to be added to the pipeline",
    type: [CreatePipelineToolDto],
  })
  @IsArray()
  @ValidateNested({ each: true })
  pipelineTools: CreatePipelineToolDto[];
}
