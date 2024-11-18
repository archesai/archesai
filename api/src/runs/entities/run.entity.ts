import { BaseEntity } from "@/src/common/entities/base.entity";
import { ApiHideProperty, ApiProperty } from "@nestjs/swagger";
import { Run as _PrismaRun, RunStatus, RunType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import {
  IsDate,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
} from "class-validator";

export type RunModel = _PrismaRun;

@Exclude()
export class RunEntity extends BaseEntity implements RunModel {
  @ApiProperty({
    description: "The timestamp when the run completed",
    example: "2024-11-05T11:42:02.258Z",
    required: false,
    type: Date,
  })
  @Expose()
  @IsOptional()
  @IsDate()
  completedAt: Date | null;

  @ApiProperty({
    description: "The error message, if any, associated with the run",
    example: "An unexpected error occurred.",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  @IsString()
  error: null | string;

  @ApiProperty({
    description: "The name of the run",
    example: "Data Processing PipelineRun",
    required: false,
    type: String,
  })
  @Expose()
  @IsString()
  @IsOptional()
  name: null | string;

  @ApiHideProperty()
  @Exclude()
  orgname: string;

  @ApiProperty({
    description: "The pipeline ID associated with the run, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  pipelineId: null | string;

  @ApiHideProperty()
  @Exclude()
  pipelineRunId: null | string;

  @ApiHideProperty()
  @Exclude()
  pipelineStepId: null | string;

  @ApiProperty({
    default: 0,
    description: "The progress of the run as a percentage",
    example: 50.5,
  })
  @Expose()
  @IsNumber()
  progress: number;

  @ApiProperty({
    description:
      "The type of run, either an individual tool run or a pipeline run",
    enum: RunType,
    required: true,
  })
  runType: RunType;

  @ApiProperty({
    description: "The timestamp when the run started",
    example: "2024-11-05T11:42:02.258Z",
    required: false,
    type: Date,
  })
  @Expose()
  @IsOptional()
  @IsDate()
  startedAt: Date | null;

  @ApiProperty({
    default: RunStatus.QUEUED,
    description: "The status of the run",
    enum: RunStatus,
  })
  @Expose()
  @IsEnum(RunStatus)
  status: RunStatus;

  @ApiProperty({
    description: "The tool ID associated with the run, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  toolId: null | string;

  constructor(toolRun: RunModel) {
    super();
    Object.assign(this, toolRun);
  }
}
