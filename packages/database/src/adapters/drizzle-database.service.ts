import type { SQL } from 'drizzle-orm'
import type { NodePgDatabase } from 'drizzle-orm/node-postgres'
import type { PgColumn, PgTable } from 'drizzle-orm/pg-core'

import { sql } from 'drizzle-orm'
import { getTableConfig } from 'drizzle-orm/pg-core'

import type { SearchQuery } from '@archesai/core'
import type { BaseEntity, BaseInsertion } from '@archesai/domain'

import { DatabaseService } from '@archesai/core'

export class DrizzleDatabaseService<
  TEntity extends BaseEntity = BaseEntity,
  TInsert extends BaseInsertion<TEntity> = BaseInsertion<TEntity>,
  TModel = TEntity,
  TDatabase extends NodePgDatabase<Record<string, unknown>> = NodePgDatabase<
    Record<string, unknown>
  >
> extends DatabaseService<TEntity, TInsert, TModel, SQL, PgTable> {
  private readonly db: TDatabase
  constructor(db: TDatabase) {
    super()
    this.db = db
  }

  public buildWhereConditions(
    table: PgTable,
    query: SearchQuery<TEntity>
  ): SQL | undefined {
    const conditions: SQL[] = []

    if (!query.filter) {
      return undefined
    }

    // Dynamic filters
    const tableName = getTableConfig(table).name
    Object.entries(query.filter).forEach(([field, filter]) => {
      {
        Object.entries(filter as Record<string, unknown>).forEach(
          ([operator, value]) => {
            const columnRef = sql`${sql.identifier(tableName)}.${sql.identifier(field)}`
            switch (operator) {
              case 'equals':
                conditions.push(sql`${columnRef} = ${value}`)
                break
              case 'gt':
                conditions.push(sql`${columnRef} > ${value}`)
                break
              case 'gte':
                conditions.push(sql`${columnRef} >= ${value}`)
                break
              case 'in':
                if (!Array.isArray(value)) {
                  throw new Error(`Value for NOT_IN operator must be an array`)
                }
                conditions.push(
                  sql`${columnRef} IN (${sql.join(value as SQL[])}))`
                )
                break
              case 'is_not_null':
                conditions.push(sql`${columnRef} IS NOT NULL`)
                break
              case 'is_null':
                conditions.push(sql`${columnRef} IS NULL`)
                break
              case 'like':
                conditions.push(sql`${columnRef} LIKE ${value}`)
                break
              case 'lt':
                conditions.push(sql`${columnRef} < ${value}`)
                break
              case 'lte':
                conditions.push(sql`${columnRef} <= ${value}`)
                break
              case 'not_equals':
                conditions.push(sql`${columnRef} != ${value}`)
                break
              case 'not_in':
                if (!Array.isArray(value)) {
                  throw new Error(`Value for NOT_IN operator must be an array`)
                }
                conditions.push(
                  sql`${columnRef} NOT IN (${sql.join(value as SQL[])})`
                )
                break
            }
          }
        )
      }
    })

    return conditions.length > 0 ? sql.join(conditions, sql` AND `) : undefined
  }

  public async count(table: PgTable, where?: SQL): Promise<number> {
    const t = this.getTableFromName(table)
    const [count] = await this.db
      .select({ count: sql<number>`COUNT(*)` })
      .from(t)
      .where(where)
    if (!count) {
      return 0
    }
    return count.count
  }

  public async delete(table: PgTable, where?: SQL): Promise<TModel[]> {
    const t = this.getTableFromName(table)
    return this.db.delete(t).where(where).returning() as unknown as Promise<
      TModel[]
    >
  }

  public execute(query: SQL): Promise<unknown> {
    return this.db.execute(query)
  }

  public getTableFromName(table: PgTable): PgTable & {
    id: PgColumn
  } {
    return table as PgTable & {
      id: PgColumn
    }
  }

  public insert(table: PgTable, values: TInsert[]): Promise<TModel[]> {
    const t = this.getTableFromName(table)
    const query = this.db.insert(t).values(values).returning()
    return query as unknown as Promise<TModel[]>
  }

  public select(table: PgTable, where?: SQL): Promise<TModel[]> {
    const t = this.getTableFromName(table)
    return this.db.select().from(t).where(where) as unknown as Promise<TModel[]>
  }

  public update(
    table: PgTable,
    values: Partial<TInsert>,
    where?: SQL
  ): Promise<TModel[]> {
    const t = this.getTableFromName(table)
    return this.db
      .update(t)
      .set(values)
      .where(where)
      .returning() as unknown as Promise<TModel[]>
  }
}
