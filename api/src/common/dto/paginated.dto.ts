import { AggregateFieldQuery } from './search-query.dto'

export class AggregateFieldResult extends AggregateFieldQuery {
  /**
   * The value of the aggregate
   * @example 10
   */
  value: number
}

export class Metadata {
  /**
   * The number of results per page
   * @example 10
   */
  limit: number

  /**
   * The current page
   * @example 1
   */
  offset: number

  /**
   * The total number of results
   * @example 100
   */
  totalResults: number
}

export class PaginatedDto<TData> {
  /**
   * The aggregates for the paginated results
   */
  aggregates?: AggregateFieldResult[]

  /**
   * The metadata for the paginated results
   */
  metadata: Metadata

  /**
   * The paginated results
   */
  results: TData[]

  constructor(input: {
    aggregates?: AggregateFieldResult[]
    metadata: Metadata
    results: TData[]
  }) {
    this.results = input.results
    this.metadata = input.metadata
    this.aggregates = input.aggregates
  }
}
