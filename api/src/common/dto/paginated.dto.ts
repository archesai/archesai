import { ApiProperty } from "@nestjs/swagger";

import { AggregateFieldQuery } from "./search-query.dto";

export class Metadata {
  @ApiProperty({ description: "The number of results per page", example: 10 })
  limit: number;

  @ApiProperty({ description: "The current page", example: 1 })
  offset: number;

  @ApiProperty({ description: "The total number of results", example: 100 })
  totalResults: number;
}

export class AggregateFieldResult extends AggregateFieldQuery {
  @ApiProperty({
    description: "The value of the aggregate",
    example: 10,
  })
  value: number;
}

export class PaginatedDto<TData> {
  @ApiProperty({
    description: "The aggregates for the paginated results",
    isArray: true,
  })
  aggregates: AggregateFieldResult[];

  @ApiProperty({ description: "The metadata for the paginated results" })
  metadata: Metadata;

  @ApiProperty({
    description: "The paginated results",
    isArray: true,
  })
  results: TData[];

  constructor(input: {
    aggregates?: AggregateFieldResult[];
    metadata: Metadata;
    results: TData[];
  }) {
    this.results = input.results;
    this.metadata = input.metadata;
    this.aggregates = input.aggregates;
  }
}
