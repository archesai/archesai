import { TimedProcessEntity } from "@/src/common/entities/timed-process.entity";
import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { PipelineRun as _PrismaPipelineRun } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import { IsOptional } from "class-validator";

type PipelineRunModel = _PrismaPipelineRun;

@Exclude()
export class PipelineRunEntity
  extends TimedProcessEntity
  implements PipelineRunModel
{
  @ApiHideProperty()
  orgname: string;

  @ApiProperty({
    description: "The pipeline ID associated with the run, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  pipelineId: null | string;

  @ApiProperty({
    description: "The thread that the run is associated with, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  threadId: null | string;

  constructor(run: PipelineRunModel) {
    super();
    Object.assign(this, run);
  }
}
