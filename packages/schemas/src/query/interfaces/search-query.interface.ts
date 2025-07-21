import type { BaseEntity } from '#base/entities/base.entity'
import type { FilterOperation } from '#query/operator.schema'

export interface FilterCondition<TEntity extends BaseEntity> {
  field: keyof TEntity
  operator: FilterOperation
  type: 'condition'
  value:
    | (boolean | number | string)[]
    | boolean
    | number
    | string
    | { from: number | string; to: number | string }
    | { unit: 'days' | 'months' | 'weeks' | 'years'; value: number }
}

export interface FilterGroup<TEntity extends BaseEntity> {
  children: FilterNode<TEntity>[]
  operator: 'and' | 'or'
  type: 'group'
}

export type FilterNode<TEntity extends BaseEntity> =
  | FilterCondition<TEntity>
  | FilterGroup<TEntity>

export interface SearchQuery<TEntity extends BaseEntity> {
  filter?: FilterNode<TEntity>
  page?: {
    number?: number
    size?: number
  }
  sort?: {
    field: keyof TEntity
    order: 'asc' | 'desc'
  }[]
}
