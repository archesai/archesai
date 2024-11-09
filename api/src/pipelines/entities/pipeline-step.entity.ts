import { BaseEntity } from "@/src/common/entities/base.entity";
import { ToolEntity } from "@/src/tools/entities/tool.entity";
import { ApiProperty } from "@nestjs/swagger";
import {
  PipelineStep as _PrismaPipelineStep,
  Tool as _PrismaTool,
} from "@prisma/client";
import { Exclude, Expose, Transform } from "class-transformer";

type PipelineStepModel = _PrismaPipelineStep & {
  tool: _PrismaTool;
};

@Exclude()
export class PipelineStepEntity
  extends BaseEntity
  implements PipelineStepModel
{
  @ApiProperty({
    description: "The ID of the pie",
    required: false,
    type: String,
  })
  @Expose()
  dependsOnId: null | string;

  @ApiProperty({
    description: "The ID of the pipelin that this step belongs to",
    type: String,
  })
  @Expose()
  pipelineId: string;

  @ApiProperty({
    description: "The name of the tool",
    example: "My Tool",
    type: String,
  })
  @Expose()
  @Transform(({ value }) => value.name)
  tool: ToolEntity;

  @ApiProperty({
    description: "The ID of the tool that this step uses",
    type: String,
  })
  @Expose()
  toolId: string;

  constructor(pipelineStep: PipelineStepModel) {
    super();
    Object.assign(this, pipelineStep);
  }
}
