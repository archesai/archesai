import { BadRequestException } from '@nestjs/common'
import { ApiProperty, getSchemaPath } from '@nestjs/swagger'
import { Transform, Type } from 'class-transformer'
import { IsArray, IsDateString, IsEnum, IsNumber, IsOptional, IsPositive, IsString, ValidateIf } from 'class-validator'

export enum Granularity {
  DAY = 'day',
  MONTH = 'month',
  WEEK = 'week',
  YEAR = 'year'
}

export enum Operator {
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

export enum SortDirection {
  ASCENDING = 'asc',
  DESCENDING = 'desc'
}

export class AggregateFieldQuery {
  @ApiProperty({ description: 'Field to aggregate by', type: String })
  @IsString()
  field: string

  /**
   *The granularity to use for ranged aggregates
   * @example day
   */
  @IsEnum(Granularity, { always: false })
  @IsOptional()
  granularity?: Granularity

  /**
   *The type of aggregate to perform
   * @example count
   */
  @IsString()
  type: 'count' | 'sum'
}

export class FieldFieldQuery {
  @ApiProperty({ description: 'Field to filter by', type: String })
  @IsString()
  field: string

  @ApiProperty({
    description: 'Operator to use for filtering',
    enum: Operator,
    required: false
  })
  operator?: Operator = Operator.CONTAINS

  @ApiProperty({
    description: 'Value to filter for',
    oneOf: [{ type: 'string' }, { items: { type: 'string' }, type: 'array' }]
  })
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
  @IsDateString()
  @IsOptional()
  endDate?: string

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

  @ApiProperty({
    default: 10,
    description: 'The limit of the number of results returned',
    required: false
  })
  @IsNumber()
  @IsOptional()
  @IsPositive()
  limit?: number = 10

  /**
   *The offset of the returned results
   * @example 10
   */
  @IsNumber()
  @IsOptional()
  offset?: number = 0

  /**
   *The field to sort the results by
   * @example createdAt
   */
  @IsOptional()
  @IsString()
  sortBy?: string = 'createdAt'

  /**
   *The direction to sort the results by
   * @example desc
   */
  @IsEnum(SortDirection)
  @IsOptional()
  sortDirection?: SortDirection = SortDirection.DESCENDING

  /**
   *The start date to search from
   * @example 2021-01-01
   */
  @IsDateString()
  @IsOptional()
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
      throw new BadRequestException('Invalid filters format. It should be a JSON array.')
    }
  } else {
    const filters = value.map((filter: string) => JSON.parse(filter))
    return filters
  }
}
