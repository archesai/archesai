import { BaseEntity } from "@/src/common/dto/base.entity.dto";
import { ApiProperty } from "@nestjs/swagger";
import { Pipeline, PipelineTool, Tool } from "@prisma/client";
import { IsString } from "class-validator";

import { PipelineToolEntity } from "./pipeline-tool.entity";

export class PipelineEntity extends BaseEntity implements Pipeline {
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

  @ApiProperty({ type: [PipelineToolEntity] })
  pipelineTools: PipelineToolEntity[];

  constructor(
    pipeline: { pipelineTools: ({ tool: Tool } & PipelineTool)[] } & Pipeline
  ) {
    super();
    Object.assign(this, pipeline);
    this.pipelineTools = pipeline.pipelineTools.map(
      (pipelineTool) => new PipelineToolEntity(pipelineTool)
    );
  }
}
