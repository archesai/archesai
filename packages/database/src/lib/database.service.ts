import type { SQL } from 'drizzle-orm'
import type { PgTable } from 'drizzle-orm/pg-core'

import { asc, desc, sql } from 'drizzle-orm'
import { drizzle } from 'drizzle-orm/node-postgres'
import { getTableConfig } from 'drizzle-orm/pg-core'
import pg from 'pg'

import type {
  BaseEntity,
  FilterCondition,
  FilterGroup,
  FilterNode,
  SearchQuery
} from '@archesai/schemas'

import * as schema from '#schema/index'

export const createDatabaseService = (connectionString: string) => {
  const pool = new pg.Pool({ connectionString })
  const db = drizzle(pool, {
    schema
  })
  return {
    buildOrderBy<TData extends BaseEntity>(
      table: PgTable,
      query: SearchQuery<TData>
    ) {
      if (!query.sort || query.sort.length === 0) {
        return []
      }
      const tableName = getTableConfig(table).name
      return query.sort.map((sort) => {
        const columnRef = sql`${sql.identifier(tableName)}.${sql.identifier(String(sort.field))}`
        return sort.order === 'asc' ? asc(columnRef) : desc(columnRef)
      })
    },
    buildWhereConditions<TData extends BaseEntity>(
      table: PgTable,
      query: SearchQuery<TData>
    ): SQL | undefined {
      if (!query.filter) {
        return undefined
      }

      const tableName = getTableConfig(table).name

      // Recursive function to build conditions from filter tree
      const buildFilterNode = (node: FilterNode<TData>): SQL => {
        if (node.type === 'condition') {
          return buildCondition(node)
        } else {
          return buildGroup(node)
        }
      }

      // Build a single condition
      const buildCondition = (condition: FilterCondition<TData>): SQL => {
        const { field, operator, value } = condition
        const columnRef = sql`${sql.identifier(tableName)}.${sql.identifier(String(field))}`

        switch (operator) {
          case 'eq':
            return sql`${columnRef} = ${value}`

          case 'gt':
            return sql`${columnRef} > ${value}`

          case 'gte':
            return sql`${columnRef} >= ${value}`

          case 'iLike':
            return sql`${columnRef} ILIKE ${value}`

          case 'inArray':
            if (!Array.isArray(value)) {
              throw new Error(`Value for inArray operator must be an array`)
            }
            if (value.length === 0) {
              return sql`FALSE` // Empty array means no matches
            }
            return sql`${columnRef} IN (${sql.join(
              value.map((v) => sql`${v}`),
              sql`, `
            )})`

          case 'isBetween':
            if (
              typeof value !== 'object' ||
              !('from' in value) ||
              !('to' in value)
            ) {
              throw new Error(
                `Value for isBetween operator must be an object with 'from' and 'to' properties`
              )
            }
            return sql`${columnRef} BETWEEN ${value.from} AND ${value.to}`

          case 'isEmpty':
            return sql`${columnRef} IS NULL`

          case 'isNotEmpty':
            return sql`${columnRef} IS NOT NULL`

          case 'isRelativeToToday':
            if (
              typeof value !== 'object' ||
              !('value' in value) ||
              !('unit' in value)
            ) {
              throw new Error(
                `Value for isRelativeToToday operator must be an object with 'value' and 'unit' properties`
              )
            }
            return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${sql.raw(`${value.value.toString()} ${value.unit}`)}')`

          case 'lt':
            return sql`${columnRef} < ${value}`

          case 'lte':
            return sql`${columnRef} <= ${value}`

          case 'ne':
            return sql`${columnRef} != ${value}`

          case 'notILike':
            return sql`${columnRef} NOT ILIKE ${value}`

          case 'notInArray':
            if (!Array.isArray(value)) {
              throw new Error(`Value for notInArray operator must be an array`)
            }
            if (value.length === 0) {
              return sql`TRUE` // Empty array means all match
            }
            return sql`${columnRef} NOT IN (${sql.join(
              value.map((v) => sql`${v}`),
              sql`, `
            )})`

          default:
            throw new Error(`Unknown operator`)
        }
      }

      // Build a group of conditions
      const buildGroup = (group: FilterGroup<TData>): SQL => {
        if (!group.children.length || group.children.length === 0) {
          throw new Error('Filter group must have at least one child')
        }

        const childConditions = group.children.map((child) =>
          buildFilterNode(child)
        )

        if (childConditions.length === 1 && childConditions[0]) {
          return childConditions[0]
        }

        const joinOperator = group.operator === 'and' ? sql` AND ` : sql` OR `
        return sql`(${sql.join(childConditions, joinOperator)})`
      }

      return buildFilterNode(query.filter)
    },
    async count(table: PgTable, where?: SQL): Promise<number> {
      const count = await db.$count(table, where)
      return count
    },
    db,
    async delete<T extends PgTable>(table: T, where?: SQL) {
      return db.delete(table).where(where).returning()
    },
    async insert<T extends PgTable>(table: T, values: T['$inferInsert'][]) {
      return db.insert(table).values(values).returning()
    },
    async ping(): Promise<boolean> {
      try {
        await db.execute(sql`SELECT 1`)
        return true
      } catch {
        return false
      }
    },
    async select(
      table: PgTable,
      where?: SQL,
      orderBy: SQL[] = [],
      limit = 10,
      offset = 0
    ) {
      const query = await db
        .select()
        .from(table)
        .where(where)
        .orderBy(...orderBy)
        .limit(limit)
        .offset(offset)

      return query
    },
    async update<T extends PgTable>(
      table: T,
      values: T['$inferInsert'],
      where?: SQL
    ): Promise<T['$inferSelect'][]> {
      const result = await db.update(table).set(values).where(where).returning()
      return result as T['$inferSelect'][]
    }
  }
}

export type DatabaseService = ReturnType<typeof createDatabaseService>
