import { Expose } from 'class-transformer'
import { FieldAggregate } from './search-query.dto'
import { IsArray, IsNumber, ValidateNested } from 'class-validator'

export class AggregateFieldResult extends FieldAggregate {
  /**
   * The value of the aggregate
   * @example 10
   */
  @Expose()
  value: number
}

export class Metadata {
  /**
   * The number of results per page
   * @example 10
   */
  @IsNumber()
  @Expose()
  limit: number

  /**
   * The current page
   * @example 1
   */
  @IsNumber()
  @Expose()
  offset: number

  /**
   * The total number of results
   * @example 100
   */
  @IsNumber()
  @Expose()
  totalResults: number
}

export class PaginatedDto<TData> {
  /**
   * The aggregates for the paginated results
   */
  @IsArray()
  @ValidateNested({
    each: true
  })
  @Expose()
  aggregates: AggregateFieldResult[]

  /**
   * The metadata for the paginated results
   */
  @ValidateNested()
  @Expose()
  metadata: Metadata

  /**
   * The paginated results
   */
  @IsArray()
  @ValidateNested({
    each: true
  })
  @Expose()
  results: TData[]

  constructor(input: {
    aggregates: AggregateFieldResult[]
    metadata: Metadata
    results: TData[]
  }) {
    this.results = input.results
    this.metadata = input.metadata
    this.aggregates = input.aggregates
  }
}
