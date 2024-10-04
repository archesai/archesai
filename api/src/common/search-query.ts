import { ApiProperty } from "@nestjs/swagger";
import {
  IsDateString,
  IsEnum,
  IsNumber,
  IsOptional,
  IsPositive,
  IsString,
} from "class-validator";

export enum SortByField {
  CREATED = "createdAt",
}

export enum SortDirection {
  ASCENDING = "asc",
  DESCENDING = "desc",
}

export class SearchQueryDto {
  @ApiProperty({
    description: "The end date to search to",
    required: false,
  })
  @IsDateString()
  @IsOptional()
  endDate?: string;

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

  @ApiProperty({ default: "", description: "Search term", required: false })
  @IsString()
  @IsOptional()
  searchTerm?: string = "";

  @ApiProperty({
    default: SortByField.CREATED,
    description: "The field to sort the results by",
    enum: SortByField,
    required: false,
  })
  @IsEnum(SortByField)
  @IsOptional()
  sortBy? = "createdAt";

  @ApiProperty({
    default: SortDirection.DESCENDING,
    description: "The direction to sort the results by",
    enum: SortDirection,
    required: false,
  })
  @IsEnum(SortDirection)
  @IsOptional()
  sortDirection? = "desc" as SortDirection;

  @ApiProperty({
    description: "The start date to search from",
    required: false,
  })
  @IsDateString()
  @IsOptional()
  startDate?: string;
}
