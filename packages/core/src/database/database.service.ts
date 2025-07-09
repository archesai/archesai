import type { BaseEntity, BaseInsertion } from '@archesai/schemas'

import type { SearchQuery } from '#http/dto/search-query.dto'

/**
 * Abstract base class for database operations.
 *
 * @template TEntity - Domain entity shape
 * @template TInsert - Type used for insert/update operations
 * @template TModel  - Returned database row shape (defaults to TEntity)
 * @template TFilter - Custom filtering structure (e.g., SQL WHERE, Drizzle conditions, etc.)
 */
export abstract class DatabaseService<
  TEntity extends BaseEntity = BaseEntity,
  TInsert extends BaseInsertion<TEntity> = BaseInsertion<TEntity>,
  TModel = TEntity,
  TFilter = unknown,
  TTables = unknown
> {
  /**
   * Builds database filter conditions from a search query.
   */
  public abstract buildWhereConditions(
    table: TTables,
    query: SearchQuery<TEntity>
  ): TFilter | undefined

  /**
   * Returns the count of records matching the filter.
   */
  public abstract count(table: TTables, where?: TFilter): Promise<number>

  /**
   * Runs a raw database query.
   */
  public abstract delete(table: TTables, where?: TFilter): Promise<TModel[]>

  /**
   * Runs a raw database query.
   */
  public abstract execute(query: TFilter): Promise<unknown>

  /**
   * Resolves a table name to its mapped definition (e.g., a Drizzle schema).
   */
  public abstract getTableFromName(table: TTables): unknown

  /**
   * Inserts one or more records into the table.
   */
  public abstract insert(table: TTables, values: TInsert[]): Promise<TModel[]>

  /**
   * Selects records that match the filter.
   */
  public abstract select(table: TTables, where?: TFilter): Promise<TModel[]>

  /**
   * Updates records that match the filter.
   */
  public abstract update(
    table: TTables,
    values: Partial<TInsert>,
    where?: TFilter
  ): Promise<TModel[]>
}
