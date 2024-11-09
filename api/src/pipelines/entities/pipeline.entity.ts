import { BaseEntity } from "@/src/common/entities/base.entity";
import { ApiProperty } from "@nestjs/swagger";
import {
  Pipeline as _PrismaPipeline,
  PipelineStep as _PrismaPipelineStep,
  Tool as _PrismaTool,
} from "@prisma/client";
import { IsString } from "class-validator";

import { PipelineStepEntity } from "./pipeline-step.entity";

export type PipelineWithPipelineStepsModel = _PrismaPipeline & {
  pipelineSteps: (_PrismaPipelineStep & { tool: _PrismaTool })[];
};

export class PipelineEntity
  extends BaseEntity
  implements PipelineWithPipelineStepsModel
{
  @ApiProperty({
    example: "This is a sample pipeline",
    required: false,
    type: String,
  })
  @IsString()
  description: null | string;

  @ApiProperty({ example: "My Pipeline" })
  @IsString()
  name: string;

  @ApiProperty({ example: "my-organization" })
  orgname: string;

  @ApiProperty({ type: [PipelineStepEntity] })
  pipelineSteps: PipelineStepEntity[];

  constructor(pipeline: PipelineWithPipelineStepsModel) {
    super();
    Object.assign(this, pipeline);
    this.pipelineSteps = pipeline.pipelineSteps.map(
      (pipelineStep) => new PipelineStepEntity(pipelineStep)
    );
  }
}
