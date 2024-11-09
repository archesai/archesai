import { BaseEntity } from "@/src/common/entities/base.entity";
import { ApiProperty } from "@nestjs/swagger";
import { PipelineStep as _PrismaPipelineStep } from "@prisma/client";
import { Exclude } from "class-transformer";

type PipelineStepModel = _PrismaPipelineStep;

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
  dependsOnId: null | string;

  @ApiProperty({
    description: "The ID of the pipelin that this step belongs to",
    required: false,
    type: String,
  })
  pipelineId: string;

  @ApiProperty({
    description: "The ID of the tool that this step uses",
    required: false,
    type: String,
  })
  toolId: string;

  constructor(pipelineStep: PipelineStepModel) {
    super();
    Object.assign(this, pipelineStep);
  }
}
