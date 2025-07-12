import { AssertError } from '@sinclair/typebox/value'

import type { BaseEntity, TSchema } from '@archesai/schemas'

import { Value } from '@archesai/schemas'

import type { DatabaseService } from '#database/database.service'
import type { EntityFilter, SearchQuery } from '#http/dto/search-query.dto'

import { NotFoundException } from '#exceptions/http-errors'
import { Logger } from '#logging/logger'

/**
 * A base repository for handling CRUD operations on a resource.
 */
export abstract class BaseRepository<
  TEntity extends BaseEntity,
  TInsertModel = unknown,
  TSelectModel extends BaseEntity = BaseEntity,
  TTables = unknown
> {
  protected readonly entitySchema: TSchema
  protected readonly primaryKey: string = 'id'
  private readonly databaseService: DatabaseService<
    TInsertModel,
    TSelectModel,
    unknown,
    TTables
  >
  private readonly logger: Logger = new Logger(this.constructor.name)
  private readonly table: TTables

  constructor(
    databaseService: DatabaseService<
      TInsertModel,
      TSelectModel,
      unknown,
      TTables
    >,
    table: TTables,
    entitySchema: TSchema
  ) {
    this.databaseService = databaseService
    this.table = table
    this.entitySchema = entitySchema
  }

  public async create(value: TInsertModel): Promise<TEntity> {
    this.logger.debug('create', { value })
    const [model] = await this.databaseService.insert(this.table, [value])
    if (!model) {
      throw new Error('Failed to create entity')
    }
    return this.toEntity(model)
  }

  public async createMany(values: TInsertModel[]): Promise<{
    count: number
    data: TEntity[]
  }> {
    this.logger.debug('createMany', { values })
    const models = await this.databaseService.insert(this.table, values)
    return {
      count: models.length,
      data: models.map((model) => this.toEntity(model))
    }
  }

  public async delete(id: string): Promise<TEntity> {
    this.logger.debug('delete', { id })
    const query = this.buildSearchQueryPrimaryKey(id)
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const [model] = await this.databaseService.delete(
      this.table,
      whereConditions
    )
    if (!model) {
      throw new NotFoundException(`${id} not found`)
    }
    return this.toEntity(model)
  }

  public async deleteMany(
    query: SearchQuery<TSelectModel>
  ): Promise<{ count: number; data: TEntity[] }> {
    this.logger.debug('deleteMany', { query })
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const models = await this.databaseService.delete(
      this.table,
      whereConditions
    )
    return {
      count: models.length,
      data: models.map((res) => this.toEntity(res))
    }
  }

  public async findFirst(query: SearchQuery<TSelectModel>): Promise<TEntity> {
    this.logger.debug('findFirst', { query })
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      {
        ...query,
        page: {
          number: 0,
          size: 1
        }
      }
    )
    const [model] = await this.databaseService.select(
      this.table,
      whereConditions
    )
    if (!model) {
      throw new NotFoundException()
    }
    return this.toEntity(model)
  }

  public async findMany(query: SearchQuery<TSelectModel>): Promise<{
    count: number
    data: TEntity[]
  }> {
    this.logger.debug('findMany', { query })
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const models = await this.databaseService.select(
      this.table,
      whereConditions
    )
    const count = await this.databaseService.count(this.table, whereConditions)
    return {
      count: count,
      data: models.map((res) => this.toEntity(res))
    }
  }

  public async findOne(id: string): Promise<TEntity> {
    this.logger.debug('findOne', { id })
    const query = this.buildSearchQueryPrimaryKey(id)
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const [model] = await this.databaseService.select(
      this.table,
      whereConditions
    )
    if (!model) {
      throw new NotFoundException(`${id} not found`)
    }
    return this.toEntity(model)
  }

  public async update(
    id: string,
    value: Partial<TInsertModel>
  ): Promise<TEntity> {
    this.logger.debug('update', { id, value })
    const query = this.buildSearchQueryPrimaryKey(id)
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const [model] = await this.databaseService.update(
      this.table,
      value,
      whereConditions
    )
    if (!model) {
      throw new Error('Failed to update entity')
    }
    return this.toEntity(model)
  }

  public async updateMany(
    value: Partial<TInsertModel>,
    query: SearchQuery<TSelectModel>
  ): Promise<{
    count: number
    data: TEntity[]
  }> {
    this.logger.debug('updateMany', { query, value })
    const whereConditions = this.databaseService.buildWhereConditions(
      this.table,
      query
    )
    const models = await this.databaseService.update(
      this.table,
      value,
      whereConditions
    )
    return {
      count: models.length,
      data: models.map((model) => this.toEntity(model))
    }
  }

  protected toEntity(model: TSelectModel): TEntity {
    this.logger.debug('toEntity', { model })
    try {
      return Value.Parse(this.entitySchema, model)
    } catch (error) {
      if (error instanceof AssertError) {
        const errors = [...error.Errors()]
        this.logger.error('Validation error while parsing model to entity', {
          ...error,
          errors: errors,
          model
        })
        throw new Error('Validation error while parsing model to entity')
      }
      this.logger.error('Failed to parse model to entity', { error, model })
      throw new Error('Failed to parse model to entity')
    }
  }

  private buildSearchQueryPrimaryKey(value: string): SearchQuery<TSelectModel> {
    this.logger.debug('buildSearchQueryPrimaryKey', { value })
    const query: SearchQuery<TSelectModel> = {
      filter: {
        id: {
          equals: value
        }
      } as EntityFilter<TSelectModel>,
      page: {
        number: 0,
        size: 1
      },
      sort: '-createdAt'
    }
    return query
  }
}
