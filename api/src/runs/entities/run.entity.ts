import { BaseEntity } from '@/src/common/entities/base.entity'
import { ApiHideProperty } from '@nestjs/swagger'
import { Run as _PrismaRun } from '@prisma/client'
import { Expose } from 'class-transformer'

import {
  _PrismaSubItemModel,
  SubItemEntity
} from '@/src/common/entities/base-sub-item.entity'
import {
  IsDateString,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
  ValidateNested
} from 'class-validator'

export type RunModel = _PrismaRun & {
  inputs: _PrismaSubItemModel[]
  outputs: _PrismaSubItemModel[]
}

export enum RunStatusEnum {
  QUEUED = 'QUEUED',
  PROCESSING = 'PROCESSING',
  COMPLETE = 'COMPLETE',
  ERROR = 'ERROR'
}

export enum RunTypeEnum {
  TOOL_RUN = 'TOOL_RUN',
  PIPELINE_RUN = 'PIPELINE_RUN'
}

export class RunEntity extends BaseEntity implements RunModel {
  /**
   *The timestamp when the run completed
   * @example 2024-11-05T11:42:02.258Z
   */
  @IsOptional()
  @IsDateString()
  @Expose()
  completedAt: Date | null

  /**
   * The error message, if any, associated with the run
   * @example 'An unexpected error occurred.'
   */
  @IsOptional()
  @IsString()
  @Expose()
  error: null | string

  /**
   * The inputs associated with the run
   */
  @ValidateNested({ each: true })
  @Expose()
  inputs: SubItemEntity[]

  /**
   * The name of the run
   * @example 'Data Processing PipelineRun'
   */
  @IsOptional()
  @IsString()
  @Expose()
  name: null | string

  @ApiHideProperty()
  orgname: string

  /**
   * The outputs associated with the run
   */
  @ValidateNested({ each: true })
  @Expose()
  outputs: SubItemEntity[]

  /**
   * The pipeline ID associated with the run, if applicable
   * @example '123e4567-e89b-12d3-a456-426614174000'
   */
  @IsOptional()
  @IsString()
  @Expose()
  pipelineId: null | string

  @ApiHideProperty()
  pipelineRunId: null | string

  @ApiHideProperty()
  pipelineStepId: null | string

  /**
   * The progress of the run as a percentage
   * @example 50.5
   */
  @IsNumber()
  @Expose()
  progress: number

  /**
   * The type of run, either an individual tool run or a pipeline run
   * @example TOOL_RUN
   */
  @IsEnum(RunTypeEnum)
  @Expose()
  runType: RunTypeEnum

  /**
   * The timestamp when the run started
   * @example '2024-11-05T11:42:02.258Z'
   */
  @IsOptional()
  @IsDateString()
  @Expose()
  startedAt: Date | null

  /**
   * The status of the run
   * @example QUEUED
   */
  @IsEnum(RunStatusEnum)
  @Expose()
  status: RunStatusEnum

  /**
   * The tool ID associated with the run, if applicable
   * @example '123e4567-e89b-12d3-a456-426614174000'
   */
  @IsOptional()
  @IsString()
  @Expose()
  toolId: null | string

  constructor(toolRun: RunModel) {
    super()
    Object.assign(this, toolRun)
  }
}
