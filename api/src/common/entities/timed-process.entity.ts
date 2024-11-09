import { ApiProperty } from "@nestjs/swagger";
import { RunStatus } from "@prisma/client";
import { Exclude, Expose } from "class-transformer";
import {
  IsDate,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
} from "class-validator";

import { BaseEntity } from "./base.entity";

@Exclude()
export class TimedProcessEntity extends BaseEntity {
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
  })
  @Expose()
  @IsString()
  name: string;

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
}
