import type { SearchQuery } from '@archesai/core'
import type { BaseEntity } from '@archesai/schemas'

import { DatabaseService } from '@archesai/core'

type InMemoryDatabase<TEntity extends BaseEntity = BaseEntity> = Record<
  string,
  InMemoryDatabaseTable<TEntity>
>
type InMemoryDatabaseTable<TEntity extends BaseEntity> = Record<string, TEntity>

/**
 * A simple in-memory database service.
 */
export class InMemoryDatabaseService<
  TEntity extends BaseEntity
> extends DatabaseService<TEntity> {
  private readonly database: InMemoryDatabase<TEntity> = {}

  public buildWhereConditions(table: string, query: SearchQuery<TEntity>) {
    return JSON.stringify({
      query,
      table
    })
  }

  public count(table: string, where?: string): Promise<number> {
    const t = this.getTableFromName(table)
    const entities = Object.values(t)
    if (where) {
      return Promise.resolve(entities.length)
    } else {
      return Promise.resolve(entities.length)
    }
  }

  public delete(table: string, where?: string): Promise<TEntity[]> {
    const t = this.getTableFromName(table)
    const entities = Object.values(t)
    if (where) {
      return Promise.resolve(entities)
    } else {
      return Promise.resolve(entities)
    }
  }

  public execute(_query: string): Promise<TEntity[]> {
    throw new Error('Method not implemented.')
  }

  public getTableFromName(table: string) {
    const t = this.database[table]
    if (!t) {
      return (this.database[table] = {})
    }
    return t
  }

  public insert(_table: string, _values: TEntity[]): Promise<TEntity[]> {
    throw new Error('Method not implemented yet')
  }

  public select(table: string, where?: string): Promise<TEntity[]> {
    const t = this.getTableFromName(table)
    const entities = Object.values(t)
    if (where) {
      return Promise.resolve(entities)
    } else {
      return Promise.resolve(entities)
    }
  }

  public update(
    table: string,
    _values: Partial<TEntity>,
    where?: string
  ): Promise<TEntity[]> {
    const t = this.getTableFromName(table)
    const entities = Object.values(t)
    if (where) {
      return Promise.resolve(entities)
    } else {
      return Promise.resolve(entities)
    }
  }
}
