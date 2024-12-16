import { BadRequestException } from '@nestjs/common'
import { ApiProperty, getSchemaPath } from '@nestjs/swagger'
import { Transform, Type } from 'class-transformer'
import {
  IsArray,
  IsDateString,
  IsEnum,
  IsNumber,
  IsOptional,
  IsString,
  ValidateIf
} from 'class-validator'

export enum GranularityEnum {
  DAY = 'day',
  MONTH = 'month',
  WEEK = 'week',
  YEAR = 'year'
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

export class AggregateFieldQuery {
  /**
   * The field to aggregate by
   * @example createdAt
   */
  @IsString()
  field: string

  /**
   *The granularity to use for ranged aggregates
   * @example day
   */
  @IsEnum(GranularityEnum)
  granularity: GranularityEnum

  /**
   *The type of aggregate to perform
   * @example count
   */
  @IsEnum(['count', 'sum'])
  type: 'count' | 'sum'
}

export class FieldFieldQuery {
  /**
   * The field to filter by
   * @example createdAt
   */
  @IsString()
  field: string

  /**
   * The operator to use for filtering
   * @example contains
   */
  @IsEnum(OperatorEnum)
  operator: OperatorEnum

  /**
   * The value to filter by
   * @example 2021-01-01
   */
  @IsArray()
  @IsString()
  @Type(() => String) // Ensures the array elements are treated as strings
  @ValidateIf((o) => typeof o.value === 'string')
  @ValidateIf((o) => Array.isArray(o.value))
  value: string | string[]
}

export class SearchQueryDto {
  @ApiProperty({
    default: [],
    description: 'Aggregates to collect for the search results',
    isArray: true,
    items: {
      $ref: getSchemaPath(AggregateFieldQuery)
    },
    required: false,
    type: 'array'
  })
  @IsArray()
  @IsOptional()
  @Transform(({ value }) => transformValues(value))
  aggregates?: AggregateFieldQuery[] = []

  /**
   *The end date to search to
   * @example 2022-01-01
   */
  @IsOptional()
  @IsDateString()
  endDate?: Date

  @ApiProperty({
    default: [],
    description: 'Filter fields and values',
    isArray: true,
    items: {
      $ref: getSchemaPath(FieldFieldQuery)
    },
    required: false,
    type: 'array'
  })
  @IsArray()
  @IsOptional()
  @Transform(({ value }) => transformValues(value))
  filters?: FieldFieldQuery[] = []

  /**
   * The limit of the number of results returned
   * @example 10
   */
  @IsNumber()
  limit?: number = 10

  /**
   *The offset of the returned results
   * @example 10
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

const transformValues = (value: string | string[]) => {
  if (typeof value === 'string') {
    try {
      const parsed = JSON.parse(value)
      if (!Array.isArray(parsed)) {
        const filters = [parsed]
        return filters
      }
      return parsed
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (error) {
      throw new BadRequestException(
        'Invalid filters format. It should be a JSON array.'
      )
    }
  } else {
    const filters = value.map((filter: string) => JSON.parse(filter))
    return filters
  }
}
