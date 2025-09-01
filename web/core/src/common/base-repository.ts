import type { BaseEntity, SearchQuery } from '@archesai/schemas'

export interface BaseRepository<TEntity extends BaseEntity, TInsert, TSelect> {
  create(data: TInsert): Promise<TEntity>
  createMany(data: TInsert[]): Promise<{ count: number; data: TEntity[] }>
  delete(id: string): Promise<TEntity>
  deleteMany(
    query: SearchQuery<TSelect>
  ): Promise<{ count: number; data: TEntity[] }>
  findMany(
    query: SearchQuery<TSelect>
  ): Promise<{ count: number; data: TEntity[] }>
  findOne(id: string): Promise<TEntity>
  update(id: string, data: Partial<TInsert>): Promise<TEntity>
  updateMany(
    data: Partial<TInsert>,
    query: SearchQuery<TSelect>
  ): Promise<{ count: number; data: TEntity[] }>
}
