import { BaseEntity } from '@/src/common/entities/base.entity'
// import { ContentEntity } from "@/src/content/entities/content.entity";
// import { PipelineEntity } from "@/src/pipelines/entities/pipeline.entity";
// import { PipelineStepEntity } from "@/src/pipelines/entities/pipeline-step.entity";
// import { ToolEntity } from "@/src/tools/entities/tool.entity";
import { ApiHideProperty, ApiProperty } from '@nestjs/swagger'
import {
  // Content as _PrismaContent,
  // Pipeline as _PrismaPipeline,
  // PipelineStep as _PrismaPipelineStep,
  Run as _PrismaRun,
  // Tool as _PrismaTool,
  RunStatus,
  RunType
} from '@prisma/client'
import { Exclude, Expose } from 'class-transformer'
import { IsDateString, IsEnum, IsNumber, IsOptional, IsString } from 'class-validator'

import { _PrismaSubItemModel, SubItemEntity } from '../../common/entities/base-sub-item.entity'

export type RunModel = _PrismaRun & {
  inputs: _PrismaSubItemModel[]
  outputs: _PrismaSubItemModel[]
  // pipeline: _PrismaPipeline;
  // pipelineRun: _PrismaRun;
  // pipelineStep: _PrismaPipelineStep;
  // tool: _PrismaTool;
  // toolRuns: _PrismaRun[];
}

@Exclude()
export class RunEntity extends BaseEntity implements RunModel {
  /**
   *The timestamp when the run completed
   * @example 2024-11-05T11:42:02.258Z
   */
  @Expose()
  @IsDateString()
  @IsOptional()
  completedAt: Date | null

  @ApiProperty({
    description: 'The error message, if any, associated with the run',
    example: 'An unexpected error occurred.',
    required: false,
    type: String
  })
  @Expose()
  @IsOptional()
  @IsString()
  error: null | string

  @ApiProperty({
    description: 'The inputs associated with the run',
    type: [SubItemEntity]
  })
  @Expose()
  inputs: SubItemEntity[]

  @ApiProperty({
    description: 'The name of the run',
    example: 'Data Processing PipelineRun',
    required: false,
    type: String
  })
  @Expose()
  @IsOptional()
  @IsString()
  name: null | string

  @ApiHideProperty()
  orgname: string

  @ApiProperty({
    description: 'The outputs associated with the run',
    required: false,
    type: [SubItemEntity]
  })
  @Expose()
  outputs: SubItemEntity[]

  @ApiProperty({
    description: 'The pipeline ID associated with the run, if applicable',
    example: '123e4567-e89b-12d3-a456-426614174000',
    required: false,
    type: String
  })
  @Expose()
  @IsOptional()
  @IsString()
  pipelineId: null | string

  // @ApiProperty({
  //   description:
  //     "The parent pipeline run associated with the run, if applicable",
  //   required: false,
  //   type: RunEntity,
  // })
  // pipelineRun: RunEntity;

  @ApiHideProperty()
  pipelineRunId: null | string

  // @ApiProperty({
  //   description: "The pipeline step associated with the run",
  //   required: false,
  //   type: PipelineStepEntity,
  // })
  // pipelineStep: PipelineStepEntity;

  @ApiHideProperty()
  pipelineStepId: null | string

  /**
   * The progress of the run as a percentage
   * @example 50.5
   */
  @Expose()
  @IsNumber()
  progress: number = 0

  @ApiProperty({
    description: 'The type of run, either an individual tool run or a pipeline run',
    enum: RunType,
    example: RunType.TOOL_RUN
  })
  @Expose()
  @IsEnum(RunType)
  runType: RunType

  /**
   * The timestamp when the run started
   * @example '2024-11-05T11:42:02.258Z'
   *
   */
  @Expose()
  @IsDateString()
  @IsOptional()
  startedAt: Date | null

  @ApiProperty({
    description: 'The status of the run',
    enum: RunStatus,
    example: RunStatus.QUEUED
  })
  @Expose()
  @IsEnum(RunStatus)
  status: RunStatus

  @ApiProperty({
    description: 'The tool ID associated with the run, if applicable',
    example: '123e4567-e89b-12d3-a456-426614174000',
    required: false,
    type: String
  })
  @Expose()
  @IsOptional()
  @IsString()
  toolId: null | string

  // @ApiProperty({
  //   description: "The tool runs associated with the run, if applicable",
  //   required: false,
  //   type: [RunEntity],
  // })
  // toolRuns: RunEntity[];

  constructor(toolRun: RunModel) {
    super()
    Object.assign(this, toolRun)
  }
}
