import type { BaseEntity } from '@archesai/schemas'

import type { SearchQuery } from '#http/dto/search-query.dto'

/**
 * Abstract base class for database operations.
 *
 * @template TEntity - Domain entity shape
 * @template TInsertModel - Type used for insert/update operations
 * @template TSelectModel  - Returned database row shape (defaults to TEntity)
 * @template TFilter - Custom filtering structure (e.g., SQL WHERE, Drizzle conditions, etc.)
 */
export abstract class DatabaseService<
  TInsertModel = unknown,
  TSelectModel extends BaseEntity = BaseEntity,
  TFilter = unknown,
  TTable = unknown
> {
  /**
   * Builds database filter conditions from a search query.
   */
  public abstract buildWhereConditions(
    table: TTable,
    query: SearchQuery<TSelectModel>
  ): TFilter | undefined

  /**
   * Returns the count of records matching the filter.
   */
  public abstract count(table: TTable, where?: TFilter): Promise<number>

  /**
   * Runs a raw database query.
   */
  public abstract delete(
    table: TTable,
    where?: TFilter
  ): Promise<TSelectModel[]>

  /**
   * Runs a raw database query.
   */
  public abstract execute(query: TFilter): Promise<unknown>

  /**
   * Inserts one or more records into the table.
   */
  public abstract insert(
    table: TTable,
    values: TInsertModel[]
  ): Promise<TSelectModel[]>

  /**
   * Selects records that match the filter.
   */
  public abstract select(
    table: TTable,
    where?: TFilter
  ): Promise<TSelectModel[]>

  /**
   * Updates records that match the filter.
   */
  public abstract update(
    table: TTable,
    values: Partial<TInsertModel>,
    where?: TFilter
  ): Promise<TSelectModel[]>
}
