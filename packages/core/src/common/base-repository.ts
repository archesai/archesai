import type { DatabaseService, EntityFilter, SearchQuery } from '@archesai/core'
import type { BaseEntity, TSchema } from '@archesai/schemas'

import { NotFoundException } from '@archesai/core'
import { Value } from '@archesai/schemas'

export type BaseRepository<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
> = ReturnType<typeof createBaseRepository<TEntity, TInsert, TSelect>>

export function createBaseRepository<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
>(
  databaseService: DatabaseService<TInsert, TSelect>,
  table: unknown,
  entitySchema: TSchema
) {
  const toEntity = (model: TSelect): TEntity => {
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

  const buildSearchQueryPrimaryKey = (value: string): SearchQuery<TSelect> => {
    const query: SearchQuery<TSelect> = {
      filter: {
        id: {
          equals: value
        }
      } as EntityFilter<TSelect>,
      page: {
        number: 0,
        size: 1
      },
      sort: '-createdAt'
    }
    return query
  }

  return {
    async create(data: TInsert): Promise<TEntity> {
      const [model] = await databaseService.insert(table, [data])
      if (!model) {
        throw new Error('Failed to create entity')
      }
      return toEntity(model)
    },

    async createMany(
      data: TInsert[]
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
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const models = await databaseService.delete(table, whereConditions)
      return {
        count: models.length,
        data: models.map((res) => toEntity(res))
      }
    },

    async findMany(
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const models = await databaseService.select(table, whereConditions)
      const count = await databaseService.count(table, whereConditions)
      return {
        count: count,
        data: models.map((res) => toEntity(res))
      }
    },

    async findOne(id: string): Promise<TEntity> {
      const query = buildSearchQueryPrimaryKey(id)
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const [model] = await databaseService.select(table, whereConditions)
      if (!model) {
        throw new NotFoundException(`${id} not found`)
      }
      return toEntity(model)
    },

    async update(id: string, data: Partial<TInsert>): Promise<TEntity> {
      const query = buildSearchQueryPrimaryKey(id)
      const whereConditions = databaseService.buildWhereConditions(table, query)
      const [model] = await databaseService.update(table, data, whereConditions)
      if (!model) {
        throw new Error('Failed to update entity')
      }
      return toEntity(model)
    },

    async updateMany(
      data: Partial<TInsert>,
      query: SearchQuery<TSelect>
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
