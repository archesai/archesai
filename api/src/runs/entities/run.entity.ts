import { ApiProperty } from "@nestjs/swagger";
import { Run, RunStatus, RunType } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import {
  IsDate,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
} from "class-validator";

import { BaseEntity } from "../../common/base-entity.dto"; // Assuming these enums are defined in your Prisma schema

@Exclude()
export class RunEntity extends BaseEntity implements Run {
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
    example: "Data Processing Run",
  })
  @Expose()
  @IsString()
  name: string;

  @ApiProperty({
    description: "The organization name associated with the run",
    example: "my-organization",
  })
  @Expose()
  @IsString()
  orgname: string;

  @ApiProperty({
    description:
      "The parent run ID, if this run is a child run ie a tool run that is part of a pipeline run",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  parentRunId: null | string;

  @ApiProperty({
    description: "The pipeline ID associated with the run, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  pipelineId: null | string;

  @ApiProperty({
    default: 0,
    description: "The progress of the run as a percentage",
    example: 50.5,
  })
  @Expose()
  @IsNumber()
  progress: number;

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
    description: "The tool used in the run, if applicable",
    required: false,
    type: String,
  })
  @Expose()
  @IsOptional()
  toolId: null | string;

  @ApiProperty({
    description: "The type of the run",
    enum: RunType,
  })
  @Expose()
  @IsEnum(RunType)
  type: RunType;

  constructor(run: Run) {
    super();
    Object.assign(this, run);
  }
}
