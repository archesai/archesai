import type { SQL } from 'drizzle-orm'
import type { NodePgDatabase } from 'drizzle-orm/node-postgres'
import type { PgTable } from 'drizzle-orm/pg-core'

import { sql } from 'drizzle-orm'
import { getTableConfig } from 'drizzle-orm/pg-core'
import pg from 'pg'

import type { SearchQuery } from '@archesai/core'
import type { BaseEntity } from '@archesai/schemas'

import type * as schema from '#schema/index'

import { createPooledClient } from '#helpers/clients'

export const createDrizzleDatabaseService = (
  connectionString: string
): DrizzleDatabaseService => {
  const pool = new pg.Pool({ connectionString })
  const db = createPooledClient(pool)
  return new DrizzleDatabaseService(db)
}

export class DrizzleDatabaseService {
  readonly db: NodePgDatabase<typeof schema>
  constructor(db: NodePgDatabase<typeof schema>) {
    this.db = db
  }

  public buildWhereConditions(
    table: PgTable,
    query: SearchQuery<BaseEntity>
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
    const count = await this.db.$count(table, where)
    return count
  }

  public delete<T extends PgTable>(table: T, where?: SQL) {
    return this.db.delete(table).where(where).returning()
  }

  public execute(query: SQL): Promise<unknown> {
    return this.db.execute(query)
  }

  public insert<T extends PgTable>(table: T, values: T['$inferInsert'][]) {
    return this.db.insert(table).values(values).returning()
  }

  public async ping(): Promise<boolean> {
    try {
      await this.db.execute(sql`SELECT 1`)
      return true
    } catch {
      return false
    }
  }

  public select<T extends PgTable>(table: T, where?: SQL) {
    //@ts-ignore
    return this.db.select().from(table).where(where)
  }

  public update<T extends PgTable>(
    table: T,
    values: T['$inferInsert'][],
    where?: SQL
  ) {
    return this.db.update(table).set(values).where(where).returning()
  }
}
