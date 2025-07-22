import type { SortingState } from '@tanstack/react-table'

import { useNavigate, useSearch } from '@tanstack/react-router'

import type {
  BaseEntity,
  FilterCondition,
  FilterNode,
  FilterValue,
  SearchQuery
} from '@archesai/schemas'

import { SearchQuerySchema } from '@archesai/schemas'

import type { FilterOperator } from '#types/simple-data-table'

export interface FilterActions<TEntity extends BaseEntity> {
  addCondition: (condition: FilterCondition<TEntity>) => void
  addGroup: (
    operator: 'and' | 'or',
    conditions: FilterCondition<TEntity>[]
  ) => void
  addSort: (field: keyof TEntity, order: 'asc' | 'desc') => void
  removeCondition: (field: keyof TEntity) => void
  resetAll: () => void
  resetFilters: () => void
  resetPagination: () => void
  resetSorting: () => void
  setCondition: (
    field: keyof TEntity,
    operator: FilterOperator,
    value: FilterValue
  ) => void
  setFilter: (filter: FilterNode<TEntity> | undefined) => void
  setPage: (page: number) => void
  setPageSize: (size: number) => void
  setSearchQuery: (query: SearchQuery<TEntity>) => void
  setSorting: (sorting: SortingState) => void
  wrapInGroup: (operator: 'and' | 'or') => void
}

export interface FilterState<TEntity extends BaseEntity> {
  filter?: FilterNode<TEntity> | undefined
  hasFilters: boolean
  hasSorting: boolean
  isEmpty: boolean
  pageNumber: number
  pageSize: number
  searchQuery: SearchQuery<TEntity>
  sorting: SortingState
}

/**
 * Simplified filter state management that supports your complex nested filters
 * but removes unnecessary frontend complexity.
 */
export function useFilterState<
  TEntity extends BaseEntity
>(): FilterActions<TEntity> & FilterState<TEntity> {
  const search = useSearch({ strict: false }) as unknown
  const navigate = useNavigate()

  // Parse current search query
  const searchQuery = SearchQuerySchema.parse(search) as SearchQuery<TEntity>

  // Extract convenience values
  const filter = searchQuery.filter
  const pageNumber = searchQuery.page?.number ?? 1
  const pageSize = searchQuery.page?.size ?? 10
  const sorting: SortingState = (searchQuery.sort ?? []).map((s) => ({
    desc: s.order === 'desc',
    id: String(s.field)
  }))

  // Computed state
  const hasFilters = !!filter
  const hasSorting = (searchQuery.sort?.length ?? 0) > 0
  const isEmpty = !hasFilters && !hasSorting && pageNumber === 1

  // Update search params
  const updateSearch = (updates: Partial<SearchQuery<TEntity>>) => {
    void navigate({
      replace: true,
      search: (prev: SearchQuery<BaseEntity>) =>
        ({
          ...prev,
          ...updates,
          page: {
            ...prev.page,
            ...updates.page
          }
        }) as never
    })
  }

  // Direct query manipulation
  const setSearchQuery = (query: SearchQuery<TEntity>) => {
    void navigate({
      replace: true,
      search: () => query as never
    })
  }

  // Filter operations
  const setFilter = (newFilter: FilterNode<TEntity> | undefined) => {
    if (!newFilter) {
      const { filter: _, ...rest } = searchQuery
      updateSearch(rest)
      return
    }
    updateSearch({ filter: newFilter })
  }

  const addCondition = (condition: FilterCondition<TEntity>) => {
    const currentFilter = filter

    if (!currentFilter) {
      setFilter(condition)
      return
    }

    if (currentFilter.type === 'condition') {
      // Wrap existing condition and new condition in AND group
      setFilter({
        children: [currentFilter, condition],
        operator: 'and',
        type: 'group'
      })
    } else {
      // Add to existing group
      setFilter({
        ...currentFilter,
        children: [...currentFilter.children, condition]
      })
    }
  }

  const removeCondition = (field: keyof TEntity) => {
    if (!filter) return

    const removeFromNode = (
      node: FilterNode<TEntity>
    ): FilterNode<TEntity> | undefined => {
      if (node.type === 'condition') {
        return node.field === field ? undefined : node
      } else {
        const filteredChildren = node.children
          .map(removeFromNode)
          .filter((child): child is FilterNode<TEntity> => child !== undefined)

        if (filteredChildren.length === 0) return undefined
        if (filteredChildren.length === 1) return filteredChildren[0]

        return {
          ...node,
          children: filteredChildren
        }
      }
    }

    setFilter(removeFromNode(filter))
  }

  const setCondition = (
    field: keyof TEntity,
    operator: FilterOperator,
    value: FilterValue
  ) => {
    const condition: FilterCondition<TEntity> = {
      field,
      operator: operator,
      type: 'condition',
      value
    }

    // Remove existing condition for this field, then add new one
    removeCondition(field)
    addCondition(condition)
  }

  const addGroup = (
    operator: 'and' | 'or',
    conditions: FilterCondition<TEntity>[]
  ) => {
    const newGroup: FilterNode<TEntity> = {
      children: conditions,
      operator,
      type: 'group'
    }

    if (!filter) {
      setFilter(newGroup)
    } else {
      setFilter({
        children: [filter, newGroup],
        operator: 'and',
        type: 'group'
      })
    }
  }

  const wrapInGroup = (operator: 'and' | 'or') => {
    if (!filter) return

    setFilter({
      children: [filter],
      operator,
      type: 'group'
    })
  }

  // Pagination
  const setPage = (page: number) => {
    updateSearch({
      page: {
        number: page
      }
    })
  }

  const setPageSize = (size: number) => {
    updateSearch({
      page: {
        number: 1,
        size
      }
    })
  }

  // Sorting
  const setSorting = (newSorting: SortingState) => {
    updateSearch({
      sort: newSorting.map((s) => ({
        field: s.id as keyof TEntity,
        order: s.desc ? 'desc' : 'asc'
      }))
    })
    return newSorting
  }

  const addSort = (field: keyof TEntity, order: 'asc' | 'desc') => {
    const newSort = { desc: order === 'desc', id: String(field) }
    setSorting([...sorting.filter((s) => s.id !== String(field)), newSort])
  }

  // Reset operations
  const resetFilters = () => {
    const { filter: _, ...rest } = searchQuery
    setSearchQuery(rest)
  }

  const resetPagination = () => {
    updateSearch({ page: { number: 1, size: 10 } })
  }

  const resetSorting = () => {
    updateSearch({ sort: [] })
  }

  const resetAll = () => {
    setSearchQuery({})
  }

  return {
    addCondition,
    addGroup,
    addSort,
    filter,
    hasFilters,
    hasSorting,
    isEmpty,
    pageNumber,
    pageSize,
    removeCondition,
    resetAll,
    resetFilters,
    resetPagination,
    resetSorting,
    searchQuery,
    setCondition,
    setFilter,
    setPage,
    setPageSize,
    setSearchQuery,
    setSorting,
    sorting,
    wrapInGroup
  }
}
