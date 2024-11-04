import { BaseEntity } from "@/src/common/base-entity.dto";
import { ToolEntity } from "@/src/tools/entities/tool.entity";
import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { PipelineTool, Tool } from "@prisma/client";
import { Exclude, Transform } from "class-transformer";

export class PipelineToolEntity extends BaseEntity implements PipelineTool {
  @ApiProperty({ example: "depends-on-id-uuid", required: false, type: String })
  dependsOnId: null | string;

  @ApiHideProperty()
  @Exclude()
  pipelineId: string;

  @ApiProperty({
    description: "The name of the tool",
    example: "Tool Name",
    type: String,
  })
  @Transform(({ value }) => value.name)
  tool: ToolEntity;

  @ApiProperty({ example: "tool-id-uuid" })
  toolId: string;

  constructor(
    pipelineTool: {
      tool: Tool;
    } & PipelineTool
  ) {
    super();
    Object.assign(this, pipelineTool);
    this.tool = new ToolEntity(pipelineTool.tool);
  }
}
