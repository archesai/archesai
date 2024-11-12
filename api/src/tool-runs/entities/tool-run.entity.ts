import { TimedProcessEntity } from "@/src/common/entities/timed-process.entity";
import { ApiHideProperty } from "@nestjs/swagger";
import { ToolRun as _PrismaToolRun } from "@prisma/client";
import { Exclude } from "class-transformer";

export type ToolRunModel = _PrismaToolRun;

export class ToolRunEntity extends TimedProcessEntity implements ToolRunModel {
  @ApiHideProperty()
  @Exclude()
  orgname: string;

  @ApiHideProperty()
  @Exclude()
  pipelineRunId: string;

  @ApiHideProperty()
  @Exclude()
  pipelineStepId: string;

  constructor(toolRun: ToolRunModel) {
    super();
    Object.assign(this, toolRun);
  }
}
