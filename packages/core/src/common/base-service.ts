import type { BaseEntity } from '@archesai/schemas'

import type { BaseRepository } from '#common/base-repository'
import type { SearchQuery } from '#http/dto/search-query.dto'
import type { WebsocketsService } from '#websockets/websockets.service'

export type BaseService<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
> = ReturnType<typeof createBaseService<TEntity, TInsert, TSelect>>

export function createBaseService<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
>(
  repository: BaseRepository<TEntity, TInsert, TSelect>,
  websocketsService: undefined | WebsocketsService,
  emitMutationEvent: (
    entity: TEntity,
    websocketsService: WebsocketsService
  ) => void
) {
  return {
    async create(data: TInsert): Promise<TEntity> {
      const entity = await repository.create(data)
      if (websocketsService) {
        emitMutationEvent(entity, websocketsService)
      }
      return entity
    },

    async createMany(
      data: TInsert[]
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.createMany(data)
      if (websocketsService) {
        result.data.forEach((entity) => {
          emitMutationEvent(entity, websocketsService)
        })
      }
      return result
    },

    async delete(id: string): Promise<TEntity> {
      const entity = await repository.delete(id)
      if (websocketsService) {
        emitMutationEvent(entity, websocketsService)
      }
      return entity
    },

    async deleteMany(
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.deleteMany(query)
      if (websocketsService) {
        result.data.forEach((entity) => {
          emitMutationEvent(entity, websocketsService)
        })
      }
      return result
    },

    async findMany(query: SearchQuery<TSelect>): Promise<{
      count: number
      data: TEntity[]
    }> {
      return repository.findMany(query)
    },

    async findOne(id: string): Promise<TEntity> {
      const found = await repository.findOne(id)
      return found
    },

    async update(id: string, data: Partial<TInsert>): Promise<TEntity> {
      const entity = await repository.update(id, data)
      if (websocketsService) {
        emitMutationEvent(entity, websocketsService)
      }
      return entity
    },

    async updateMany(
      data: Partial<TInsert>,
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.updateMany(data, query)
      if (websocketsService) {
        result.data.forEach((entity) => {
          emitMutationEvent(entity, websocketsService)
        })
      }
      return result
    }
  }
}
