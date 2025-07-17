import type { PgTable } from 'drizzle-orm/pg-core'

import type { BaseEntity, SearchQuery, TSchema } from '@archesai/schemas'

import { NotFoundException } from '@archesai/core'
import { Value } from '@archesai/schemas'

import type { DrizzleDatabaseService } from '#adapters/drizzle-database.service'

export type BaseRepository<
  TEntity extends BaseEntity,
  TModel extends BaseEntity
> = ReturnType<typeof createBaseRepository<TEntity, TModel>>

export function createBaseRepository<
  TEntity extends BaseEntity,
  TModel extends BaseEntity
>(
  databaseService: DrizzleDatabaseService,
  table: PgTable,
  entitySchema: TSchema
) {
  const toEntity = (model: (typeof table)['$inferSelect']): TEntity => {
    return Value.Parse(entitySchema, model)

    // try {
    //   return Value.Parse(this.entitySchema, model)
    // } catch (error) {
    //   if (error instanceof AssertError) {
    //     const errors = [...error.Errors()]
    //     this.logger.error('Validation error while parsing model to entity', {
    //       ...error,
    //       errors: errors,
    //       model
    //     })
    //     throw new Error('Validation error while parsing model to entity')
    //   }
    //   this.logger.error('Failed to parse model to entity', { error, model })
    //   throw new Error('Failed to parse model to entity')
    // }
  }

  const buildSearchQueryPrimaryKey = (value: string): SearchQuery<TModel> => {
    const query: SearchQuery<TModel> = {
      filter: {
        field: 'id',
        operator: 'eq',
        type: 'condition',
        value: value
      },
      page: {
        number: 0,
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
    async create(data: (typeof table)['$inferInsert']): Promise<TEntity> {
      const [model] = await databaseService.insert(table, [data])
      if (!model) {
        throw new Error('Failed to create entity')
      }
      return toEntity(model)
    },

    async createMany(
      data: (typeof table)['$inferInsert'][]
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
      query: SearchQuery<TModel>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const models = await databaseService.delete(table, whereConditions)
      return {
        count: models.length,
        data: models.map((model) => toEntity(model))
      }
    },

    async findMany(
      query: SearchQuery<TModel>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const orderBy = databaseService.buildOrderBy(table, query)
      const models = await databaseService
        .select(table, whereConditions, orderBy)
        .limit(query.page?.size ?? 10)
        .offset(query.page?.number ?? 0)
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

    async update(
      id: string,
      data: (typeof table)['$inferInsert']
    ): Promise<TEntity> {
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
      data: (typeof table)['$inferInsert'][],
      query: SearchQuery<TModel>
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
