import type { SQL } from 'drizzle-orm'
import type { NodePgDatabase } from 'drizzle-orm/node-postgres'
import type { PgTable } from 'drizzle-orm/pg-core'

import { asc, desc, sql } from 'drizzle-orm'
import { getTableConfig } from 'drizzle-orm/pg-core'
import pg from 'pg'

import type {
  BaseEntity,
  FilterCondition,
  FilterGroup,
  FilterNode,
  SearchQuery
} from '@archesai/schemas'

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

  public buildOrderBy<TData extends BaseEntity>(
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
  }

  public buildWhereConditions<TData extends BaseEntity>(
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
            value === null ||
            !('from' in value) ||
            !('to' in value)
          ) {
            throw new Error(
              `Value for isBetween operator must be an object with 'from' and 'to' properties`
            )
          }
          const rangeValue = value as {
            from: number | string
            to: number | string
          }
          return sql`${columnRef} BETWEEN ${rangeValue.from} AND ${rangeValue.to}`

        case 'isEmpty':
          return sql`${columnRef} IS NULL`

        case 'isNotEmpty':
          return sql`${columnRef} IS NOT NULL`

        case 'isRelativeToToday':
          if (
            typeof value !== 'object' ||
            value === null ||
            !('value' in value) ||
            !('unit' in value)
          ) {
            throw new Error(
              `Value for isRelativeToToday operator must be an object with 'value' and 'unit' properties`
            )
          }
          const relativeValue = value as {
            unit: 'days' | 'months' | 'weeks' | 'years'
            value: number
          }

          // Build the interval string
          const intervalStr = `${relativeValue.value} ${relativeValue.unit}`

          // Calculate the date relative to today
          return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${sql.raw(intervalStr)}')`

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
          throw new Error(`Unknown operator: ${operator}`)
      }
    }

    // Build a group of conditions
    const buildGroup = (group: FilterGroup<TData>): SQL => {
      if (!group.children || group.children.length === 0) {
        throw new Error('Filter group must have at least one child')
      }

      const childConditions = group.children.map((child) =>
        buildFilterNode(child)
      )

      if (childConditions.length === 1) {
        return childConditions[0]!
      }

      const joinOperator = group.operator === 'and' ? sql` AND ` : sql` OR `
      return sql`(${sql.join(childConditions, joinOperator)})`
    }

    try {
      return buildFilterNode(query.filter)
    } catch (error) {
      console.error('Error building WHERE conditions:', error)
      throw error
    }
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

  public select<T extends PgTable>(table: T, where?: SQL, orderBy?: any[]) {
    //@ts-ignore
    let query = this.db.select().from(table)

    if (where) {
      //@ts-ignore
      query = query.where(where)
    }

    if (orderBy && orderBy.length > 0) {
      //@ts-ignore
      query = query.orderBy(...orderBy)
    }

    return query
  }

  public update<T extends PgTable>(
    table: T,
    values: T['$inferInsert'][],
    where?: SQL
  ) {
    return this.db.update(table).set(values).where(where).returning()
  }

  //   // Helper function to safely handle array operations
  //   private buildArrayCondition(
  //     columnRef: SQL,
  //     operator: 'inArray' | 'notInArray',
  //     value: unknown
  //   ): SQL {
  //     if (!Array.isArray(value)) {
  //       throw new Error(`Value for ${operator} operator must be an array`)
  //     }

  //     if (value.length === 0) {
  //       // Handle empty arrays appropriately
  //       return operator === 'inArray' ? sql`FALSE` : sql`TRUE`
  //     }

  //     // Convert values to SQL literals
  //     const sqlValues = value.map((v) => sql`${v}`)
  //     const joinedValues = sql.join(sqlValues, sql`, `)

  //     return operator === 'inArray' ?
  //         sql`${columnRef} IN (${joinedValues})`
  //       : sql`${columnRef} NOT IN (${joinedValues})`
  //   }

  //   // Additional helper function for complex date operations
  //   private buildRelativeDateCondition(
  //     columnRef: SQL,
  //     value: { unit: 'days' | 'months' | 'weeks' | 'years'; value: number }
  //   ): SQL {
  //     const { unit, value: amount } = value

  //     switch (unit) {
  //       case 'days':
  //         return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${amount} days')`
  //       case 'months':
  //         return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${amount} months')`
  //       case 'weeks':
  //         return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${amount} weeks')`
  //       case 'years':
  //         return sql`${columnRef} >= (CURRENT_DATE - INTERVAL '${amount} years')`
  //       default:
  //         throw new Error(`Unknown time unit: ${unit}`)
  //     }
  //   }
  // }
}
