import type { PgTable } from 'drizzle-orm/pg-core'

import type { BaseRepository } from '@archesai/core'
import type { BaseEntity, SearchQuery } from '@archesai/schemas'

import { NotFoundException } from '@archesai/core'

import type { DatabaseService } from '#lib/database.service'

export function createBaseRepository<
  TEntity extends BaseEntity,
  TTable extends PgTable = PgTable
>(
  databaseService: DatabaseService,
  table: TTable,
  entitySchema: {
    parse: (data: unknown) => TEntity
  }
): BaseRepository<TEntity, TTable['$inferInsert'], TTable['$inferSelect']> {
  const toEntity = (model: TTable['$inferSelect']): TEntity => {
    return entitySchema.parse(model)
  }

  const buildSearchQueryPrimaryKey = (
    value: string
  ): SearchQuery<TTable['$inferSelect']> => {
    const query: SearchQuery<TTable['$inferSelect']> = {
      filter: {
        field: 'id',
        operator: 'eq',
        type: 'condition',
        value: value
      },
      page: {
        number: 1,
        size: 1
      },
      sort: [
        {
          field: 'createdAt',
          order: 'desc'
        }
      ]
    }
    return query
  }

  return {
    async create(data: TTable['$inferInsert']): Promise<TEntity> {
      const [model] = await databaseService.insert(table, [data])
      if (!model) {
        throw new Error('Failed to create entity')
      }
      return toEntity(model)
    },

    async createMany(
      data: TTable['$inferInsert'][]
    ): Promise<{ count: number; data: TEntity[] }> {
      const models = await databaseService.insert(table, data)
      return {
        count: models.length,
        data: models.map((model) => toEntity(model))
      }
    },

    async delete(id: string): Promise<TEntity> {
      const query = buildSearchQueryPrimaryKey(id)
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const [model] = await databaseService.delete(table, whereConditions)
      if (!model) {
        throw new NotFoundException(`${id} not found`)
      }
      return toEntity(model)
    },

    async deleteMany(
      query: SearchQuery<TTable['$inferSelect']>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const models = await databaseService.delete(table, whereConditions)
      return {
        count: models.length,
        data: models.map((model) => toEntity(model))
      }
    },

    async findMany(
      query: SearchQuery<TTable['$inferSelect']>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const orderBy = databaseService.buildOrderBy(table, query)
      const models = await databaseService.select(
        table,
        whereConditions,
        orderBy,
        query.page?.size,
        (query.page?.number ?? 1) - 1
      )
      const count = await databaseService.count(table, whereConditions)
      return {
        count: count,
        data: models.map((model) => toEntity(model))
      }
    },

    async findOne(id: string): Promise<TEntity> {
      const query = buildSearchQueryPrimaryKey(id)
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const orderBy = databaseService.buildOrderBy(table, query)
      const [model] = await databaseService.select(
        table,
        whereConditions,
        orderBy
      )
      if (!model) {
        throw new NotFoundException(`${id} not found`)
      }
      return toEntity(model)
    },

    async update(id: string, data: TTable['$inferInsert']): Promise<TEntity> {
      const query = buildSearchQueryPrimaryKey(id)
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const [model] = await databaseService.update(
        table,
        [data],
        whereConditions
      )
      if (!model) {
        throw new Error('Failed to update entity')
      }
      return toEntity(model)
    },

    async updateMany(
      data: TTable['$inferInsert'],
      query: SearchQuery<TTable['$inferSelect']>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const models = await databaseService.update(table, data, whereConditions)
      return {
        count: models.length,
        data: models.map((model) => toEntity(model))
      }
    }
  }
}
