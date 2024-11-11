import { BadRequestException } from "@nestjs/common";
import { ApiProperty, getSchemaPath } from "@nestjs/swagger";
import { Transform } from "class-transformer";
import {
  IsArray,
  IsDateString,
  IsEnum,
  IsNumber,
  IsOptional,
  IsPositive,
  IsString,
} from "class-validator";

export enum Granularity {
  DAY = "day",
  MONTH = "month",
  WEEK = "week",
  YEAR = "year",
}

export enum SortDirection {
  ASCENDING = "asc",
  DESCENDING = "desc",
}

export enum Operator {
  CONTAINS = "contains",
  ENDS_WITH = "endsWith",
  EQUALS = "equals",
  EVERY = "every",
  NONE = "none",
  NOT = "not",
  SOME = "some",
  STARTS_WITH = "startsWith",
}

export class FieldFieldQuery {
  @ApiProperty({ description: "Field to filter by", type: String })
  @IsString()
  field: string;

  @ApiProperty({
    description: "Operator to use for filtering",
    enum: Operator,
    required: false,
  })
  operator?: Operator = Operator.CONTAINS;

  @ApiProperty({ description: "Value to filter for", type: String })
  @IsString()
  value: string;
}

export class AggregateFieldQuery {
  @ApiProperty({ description: "Field to aggregate by", type: String })
  @IsString()
  field: string;

  @ApiProperty({
    default: undefined,
    description: "The granularity to use for ranged aggregates",
    enum: Granularity,
    required: false,
  })
  @IsOptional()
  @IsEnum(Granularity, { always: false })
  granularity?: Granularity;

  @ApiProperty({
    description: "Type of aggregate to perform",
    enum: ["count", "sum"],
  })
  @IsString()
  type: "count" | "sum";
}

export class SearchQueryDto {
  @ApiProperty({
    default: [],
    description: "Aggregates to collect for the search results",
    isArray: true,
    items: {
      $ref: getSchemaPath(AggregateFieldQuery),
    },
    required: false,
    type: "array",
  })
  @IsOptional()
  @IsArray()
  @Transform(({ value }) => transformValues(value))
  aggregates?: AggregateFieldQuery[] = [];

  @ApiProperty({
    description: "The end date to search to",
    required: false,
  })
  @IsDateString()
  @IsOptional()
  endDate?: string;

  @ApiProperty({
    default: [],
    description: "Filter fields and values",
    isArray: true,
    items: {
      $ref: getSchemaPath(FieldFieldQuery),
    },
    required: false,
    type: "array",
  })
  @IsOptional()
  @IsArray()
  @Transform(({ value }) => transformValues(value))
  filters?: FieldFieldQuery[] = [];

  @ApiProperty({
    default: 10,
    description: "The limit of the number of results returned",
    required: false,
  })
  @IsOptional()
  @IsPositive()
  @IsNumber()
  limit?: number = 10;

  @ApiProperty({
    default: 0,
    description: "The offset of the returned results",
    required: false,
  })
  @IsOptional()
  @IsNumber()
  offset?: number = 0;

  @ApiProperty({
    default: "createdAt",
    description: "The field to sort the results by",
    required: false,
  })
  @IsString()
  @IsOptional()
  sortBy?: string = "createdAt";

  @ApiProperty({
    default: SortDirection.DESCENDING,
    description: "The direction to sort the results by",
    enum: SortDirection,
    required: false,
  })
  @IsEnum(SortDirection)
  @IsOptional()
  sortDirection?: SortDirection = SortDirection.DESCENDING;

  @ApiProperty({
    description: "The start date to search from",
    required: false,
  })
  @IsDateString()
  @IsOptional()
  startDate?: string;
}

const transformValues = (value: string | string[]) => {
  if (typeof value === "string") {
    try {
      const parsed = JSON.parse(value);
      if (!Array.isArray(parsed)) {
        const filters = [parsed];
        return filters;
      }
      return parsed;
    } catch (error) {
      throw new BadRequestException(
        "Invalid filters format. It should be a JSON array."
      );
    }
  } else {
    const filters = value.map((filter: string) => JSON.parse(filter));
    return filters;
  }
};
