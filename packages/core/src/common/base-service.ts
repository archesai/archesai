import type { BaseEntity, SearchQuery } from '@archesai/schemas'

export type BaseService<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
> = ReturnType<typeof createBaseService<TEntity, TInsert, TSelect>>

export interface IBaseRepository<
  TEntity extends BaseEntity,
  TInsert,
  TModel extends BaseEntity
> {
  create(data: TInsert): Promise<TEntity>
  createMany(data: TInsert[]): Promise<{ count: number; data: TEntity[] }>
  delete(id: string): Promise<TEntity>
  deleteMany(
    query: SearchQuery<TModel>
  ): Promise<{ count: number; data: TEntity[] }>
  findMany(
    query: SearchQuery<TModel>
  ): Promise<{ count: number; data: TEntity[] }>
  findOne(id: string): Promise<TEntity>
  update(id: string, data: TInsert): Promise<TEntity>
  updateMany(
    data: TInsert,
    query: SearchQuery<TModel>
  ): Promise<{ count: number; data: TEntity[] }>
}

export function createBaseService<
  TEntity extends BaseEntity,
  TInsert,
  TSelect extends BaseEntity
>(
  repository: IBaseRepository<TEntity, TInsert, TSelect>,
  emitMutationEvent: (entity: TEntity) => void
) {
  return {
    async create(data: TInsert): Promise<TEntity> {
      const entity = await repository.create(data)
      emitMutationEvent(entity)
      return entity
    },

    async createMany(
      data: TInsert[]
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.createMany(data)
      result.data.forEach((entity) => {
        emitMutationEvent(entity)
      })
      return result
    },

    async delete(id: string): Promise<TEntity> {
      const entity = await repository.delete(id)
      emitMutationEvent(entity)
      return entity
    },

    async deleteMany(
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.deleteMany(query)
      result.data.forEach((entity) => {
        emitMutationEvent(entity)
      })
      return result
    },

    async findMany(query: SearchQuery<TSelect>): Promise<{
      count: number
      data: TEntity[]
    }> {
      return repository.findMany(query)
    },

    async findOne(id: string): Promise<TEntity> {
      return repository.findOne(id)
    },

    async update(id: string, data: TInsert): Promise<TEntity> {
      const entity = await repository.update(id, data)
      emitMutationEvent(entity)
      return entity
    },

    async updateMany(
      data: TInsert,
      query: SearchQuery<TSelect>
    ): Promise<{ count: number; data: TEntity[] }> {
      const result = await repository.updateMany(data, query)
      result.data.forEach((entity) => {
        emitMutationEvent(entity)
      })
      return result
    }
  }
}
