import { BaseEntity } from "@/src/common/dto/base.entity.dto";
import { ApiProperty } from "@nestjs/swagger";
import { Pipeline as _PrismaPipeline } from "@prisma/client";
import { IsString } from "class-validator";

import { PipelineToolEntity, PipelineToolModel } from "./pipeline-tool.entity";

export type PipelineWithPipelineToolsModel = _PrismaPipeline & {
  pipelineTools: PipelineToolModel[];
};

export class PipelineEntity
  extends BaseEntity
  implements PipelineWithPipelineToolsModel
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

  @ApiProperty({ type: [PipelineToolEntity] })
  pipelineTools: PipelineToolEntity[];

  constructor(pipeline: PipelineWithPipelineToolsModel) {
    super();
    Object.assign(this, pipeline);
    this.pipelineTools = pipeline.pipelineTools.map(
      (pipelineTool) => new PipelineToolEntity(pipelineTool)
    );
  }
}
