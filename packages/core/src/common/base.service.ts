import type { BaseEntity, BaseInsertion } from '@archesai/schemas'

import type { BaseRepository } from '#common/base.repository'
import type { SearchQuery } from '#http/dto/search-query.dto'

import { Logger } from '#logging/logger'

/**
 * A base service for handling CRUD operations on a resource.
 */
export abstract class BaseService<
  TEntity extends BaseEntity,
  TInsert extends BaseInsertion<TEntity> = BaseInsertion<TEntity>
> {
  protected readonly logger: Logger = new Logger(this.constructor.name)
  protected readonly repository: BaseRepository<TEntity>

  constructor(repository: BaseRepository<TEntity>) {
    this.repository = repository
  }

  public async create(value: TInsert): Promise<TEntity> {
    this.logger.debug('create', { value })
    const entity = await this.repository.create(value)
    this.emitMutationEvent(entity)
    return entity
  }

  public async createMany(
    values: TInsert[]
  ): Promise<{ count: number; data: TEntity[] }> {
    this.logger.debug('createMany', { values })
    const { count, data: created } = await this.repository.createMany(values)
    if (count) {
      throw new Error('Failed to create entities')
    }
    if (created.length === 0 || !created[0]) {
      throw new Error('No entities created')
    }
    this.emitMutationEvent(created[0])
    return {
      count: created.length,
      data: created
    }
  }

  public async delete(id: string): Promise<TEntity> {
    this.logger.debug('delete', { id })
    const entity = await this.repository.delete(id)
    this.emitMutationEvent(entity)
    return entity
  }

  public async deleteMany(query: SearchQuery<TEntity>): Promise<{
    count: number
    data: TEntity[]
  }> {
    this.logger.debug('findMany', { query })
    return this.repository.deleteMany(query)
  }

  public async findMany(query: SearchQuery<TEntity>): Promise<{
    count: number
    data: TEntity[]
  }> {
    this.logger.debug('findMany', { query })
    return this.repository.findMany(query)
  }

  public async findOne(id: string): Promise<TEntity> {
    this.logger.debug('findOne', { id })
    const found = await this.repository.findOne(id)
    return found
  }

  public async update(id: string, data: Partial<TInsert>): Promise<TEntity> {
    this.logger.debug('update', { data, id })
    const entity = await this.repository.update(id, data)
    this.emitMutationEvent(entity)
    return entity
  }

  public async updateMany(
    value: Partial<TInsert>,
    query: SearchQuery<TEntity>
  ): Promise<{ count: number; data: TEntity[] }> {
    this.logger.debug('findMany', { query })
    return this.repository.updateMany(value, query)
  }

  // async updateRelationship(
  //   id: string,
  //   relationshipKey: keyof TEntity['relationships'],
  //   relatedIds: string[],
  //   action: 'add' | 'remove' | 'replace'
  // ): Promise<TEntity> {
  //   this.logger.debug('updateRelationship', {
  //     action,
  //     id,
  //     relatedIds,
  //     relationshipKey
  //   })

  //   const entity = await this.repository.findOne(id)

  //   let updatedRelationships = entity[relationshipKey]

  //   switch (action) {
  //     case 'add':
  //       updatedRelationships = [
  //         ...updatedRelationships,
  //         ...relatedIds.map((relatedId) => ({ id: relatedId }))
  //       ]
  //       break

  //     case 'remove':
  //       updatedRelationships = updatedRelationships.filter(
  //         (rel) => !relatedIds.includes(rel.id)
  //       )
  //       break

  //     case 'replace':
  //       updatedRelationships = relatedIds.map(
  //         (relatedId) => ({ id: relatedId }) as BaseEntity
  //       )
  //       break
  //   }

  //   const updatedEntity = await this.repository.update(id, {
  //     [relationshipKey]: updatedRelationships
  //   })

  //   this.emitMutationEvent(updatedEntity)
  //   return updatedEntity
  // }

  protected abstract emitMutationEvent(entity: TEntity): void
}
