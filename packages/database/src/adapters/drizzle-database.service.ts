import type { SQL } from 'drizzle-orm'
import type { NodePgDatabase } from 'drizzle-orm/node-postgres'
import type { PgTable } from 'drizzle-orm/pg-core'

import { sql } from 'drizzle-orm'
import { getTableConfig } from 'drizzle-orm/pg-core'

import type { SearchQuery } from '@archesai/core'
import type { BaseEntity } from '@archesai/schemas'

import { DatabaseService } from '@archesai/core'

import type * as schema from '#schema/index'

export class DrizzleDatabaseService<
  TEntity extends BaseEntity = BaseEntity,
  TInsertModel = unknown,
  TSelectModel extends BaseEntity = TEntity
> extends DatabaseService<TInsertModel, TSelectModel, SQL, PgTable> {
  private readonly db: NodePgDatabase<typeof schema>
  constructor(db: NodePgDatabase<typeof schema>) {
    super()
    this.db = db
  }

  public buildWhereConditions(
    table: PgTable,
    query: SearchQuery<TSelectModel>
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
    const [count] = await this.db
      .select({ count: sql<number>`COUNT(*)` })
      .from(table)
      .where(where)
    if (!count) {
      return 0
    }
    return count.count
  }

  public delete(
    table: keyof NodePgDatabase<typeof schema>['_']['schema'],
    where?: SQL
  ): Promise<TSelectModel[]> {
    return this.db.delete(table).where(where).returning()
  }

  public execute(query: SQL): Promise<unknown> {
    return this.db.execute(query)
  }

  public insert(
    table: keyof NodePgDatabase<typeof schema>['_']['schema'],
    values: TInsertModel[]
  ): Promise<TSelectModel[]> {
    return this.db.insert(table).values(values).returning()
  }

  public select(
    table: keyof NodePgDatabase<typeof schema>['_']['schema'],
    where?: SQL
  ): Promise<TSelectModel[]> {
    return this.db.select().from(table).where(where)
  }

  public update(
    table: keyof NodePgDatabase<typeof schema>['_']['schema'],
    values: Partial<TInsertModel>,
    where?: SQL
  ): Promise<TSelectModel[]> {
    return this.db.update(table).set(values).where(where).returning()
  }
}
