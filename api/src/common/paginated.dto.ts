import { ApiProperty } from "@nestjs/swagger";

export class Metadata {
  @ApiProperty({ description: "The number of results per page", example: 10 })
  limit: number;

  @ApiProperty({ description: "The current page", example: 1 })
  offset: number;

  @ApiProperty({ description: "The total number of results", example: 100 })
  totalResults: number;
}

export class _PaginatedDto {
  @ApiProperty({ description: "The metadata for the paginated results" })
  metadata: Metadata;
}

export class PaginatedDto<TData, TAggregates = any> extends _PaginatedDto {
  @ApiProperty({
    description: "The aggregates for the paginated results",
  })
  aggregates: TAggregates;

  @ApiProperty({
    description: "The paginated results",
    isArray: true,
  })
  results: TData[];

  constructor(input: {
    aggregates?: TAggregates;
    metadata: Metadata;
    results: TData[];
  }) {
    super();
    this.results = input.results;
    this.metadata = input.metadata;
    this.aggregates = input.aggregates;
  }
}
