import { BadRequestException } from '@nestjs/common'
import { Expose, plainToInstance, Transform, Type } from 'class-transformer'
import {
  IsDateString,
  IsDefined,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
  ValidateNested
} from 'class-validator'

export enum GranularityEnum {
  DAY = 'day',
  MONTH = 'month',
  WEEK = 'week',
  YEAR = 'year'
}

export enum AggregateTypeEnum {
  COUNT = 'count',
  SUM = 'sum'
}

export enum OperatorEnum {
  CONTAINS = 'contains',
  ENDS_WITH = 'endsWith',
  EQUALS = 'equals',
  EVERY = 'every',
  IN = 'in',
  NONE = 'none',
  NOT = 'not',
  SOME = 'some',
  STARTS_WITH = 'startsWith'
}

export enum SortDirectionEnum {
  ASCENDING = 'asc',
  DESCENDING = 'desc'
}

export class FieldAggregate {
  /**
   * The field to aggregate by
   * @example createdAt
   */
  @IsString()
  @Expose()
  field: string

  /**
   *The granularity to use for ranged aggregates
   * @example day
   */
  @IsEnum(GranularityEnum)
  @Expose()
  granularity: GranularityEnum

  /**
   *The type of aggregate to perform
   * @example count
   */
  @IsEnum(AggregateTypeEnum)
  @Expose()
  type: AggregateTypeEnum
}

export class FieldFilter {
  /**
   * The field to filter by
   * @example createdAt
   */
  @IsString()
  @Expose()
  field: string

  /**
   * The operator to use for filtering
   * @example contains
   */
  @IsEnum(OperatorEnum)
  @Expose()
  operator: OperatorEnum

  /**
   * The value to filter by
   * @example 2021-01-01
   */
  @IsDefined()
  @Expose()
  value: string | string[]
}

export class SearchQueryDto {
  /**
   * Aggregates to collect for the search results
   */
  @IsOptional()
  @ValidateNested({ each: true })
  @Type(() => FieldAggregate)
  @Transform(({ value }) => transformValues(value, FieldAggregate))
  aggregates?: FieldAggregate[] = []

  /**
   * Filters to apply to the search results
   */
  @IsOptional()
  @ValidateNested({ each: true })
  @Type(() => FieldFilter)
  @Transform(({ value }) => transformValues(value, FieldFilter))
  filters?: FieldFilter[] = []

  /**
   *The end date to search to
   * @example 2022-01-01
   */
  @IsOptional()
  @IsDateString()
  endDate?: Date

  /**
   * The limit of the number of results returned
   * @example 10
   */
  @IsNumber()
  limit?: number = 10

  /**
   *The offset of the returned results
   * @example 0
   */
  @IsNumber()
  offset?: number = 0

  /**
   *The field to sort the results by
   * @example createdAt
   */
  @IsString()
  sortBy?: string = 'createdAt'

  /**
   *The direction to sort the results by
   * @example desc
   */
  @IsEnum(SortDirectionEnum)
  sortDirection?: SortDirectionEnum = SortDirectionEnum.DESCENDING

  /**
   *The start date to search from
   * @example 2021-01-01
   */
  @IsOptional()
  @IsDateString()
  startDate?: Date
}

export const transformValues = (value: any, cls: any): any[] => {
  if (!value) {
    return []
  } else if (Array.isArray(value)) {
    return value.map((v) => plainToInstance(cls, v))
  }
  try {
    const val = JSON.parse(value)
    return Array.isArray(val)
      ? val.map((v) => plainToInstance(cls, v))
      : [plainToInstance(cls, val)]
  } catch (error) {
    throw new BadRequestException(error)
  }
}
