import type { SortingState } from '@tanstack/react-table'

import { useCallback, useMemo } from 'react'
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

  // Group operations
  addGroup: (
    operator: 'and' | 'or',
    conditions: FilterCondition<TEntity>[]
  ) => void
  addSort: (field: keyof TEntity, order: 'asc' | 'desc') => void
  removeCondition: (field: keyof TEntity) => void
  resetAll: () => void

  // Reset operations
  resetFilters: () => void
  resetPagination: () => void

  resetSorting: () => void
  setCondition: (
    field: keyof TEntity,
    operator: FilterOperator,
    value: FilterValue
  ) => void

  // Filter operations
  setFilter: (filter: FilterNode<TEntity> | undefined) => void
  // Pagination
  setPage: (page: number) => void

  setPageSize: (size: number) => void
  // Direct query manipulation
  setSearchQuery: (query: SearchQuery<TEntity>) => void
  // Sorting
  setSorting: (sorting: SortingState) => void
  wrapInGroup: (operator: 'and' | 'or') => void
}

export interface FilterState<TEntity extends BaseEntity> {
  // Convenience accessors
  filter?: FilterNode<TEntity> | undefined

  // Computed state
  hasFilters: boolean
  hasSorting: boolean
  isEmpty: boolean
  pageNumber: number

  pageSize: number
  // Raw search query for backend
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
  const searchQuery = useMemo<SearchQuery<TEntity>>(() => {
    const validated = SearchQuerySchema.parse(search)
    return validated as SearchQuery<TEntity>
  }, [search])

  // Extract convenience values
  const filter = searchQuery.filter
  const pageNumber = searchQuery.page?.number ?? 1
  const pageSize = searchQuery.page?.size ?? 10
  const sorting: SortingState = useMemo(
    () =>
      (searchQuery.sort ?? []).map((s) => ({
        desc: s.order === 'desc',
        id: String(s.field)
      })),
    [searchQuery.sort]
  )

  // Computed state
  const hasFilters = !!filter
  const hasSorting = (searchQuery.sort?.length ?? 0) > 0
  const isEmpty = !hasFilters && !hasSorting && pageNumber === 1

  // Update search params
  const updateSearch = useCallback(
    (updates: Partial<SearchQuery<TEntity>>) => {
      void navigate({
        replace: true,
        search: (prev: SearchQuery<BaseEntity>) =>
          ({
            ...prev,
            ...updates
          }) as never
      })
    },
    [navigate]
  )

  // Direct query manipulation
  const setSearchQuery = useCallback(
    (query: SearchQuery<TEntity>) => {
      void navigate({
        replace: true,
        search: () => query as never
      })
    },
    [navigate]
  )

  // Filter operations
  const setFilter = useCallback(
    (newFilter: FilterNode<TEntity> | undefined) => {
      if (!newFilter) {
        const { filter: _, ...rest } = searchQuery
        updateSearch(rest)
        return
      }
      updateSearch({ filter: newFilter })
    },
    [updateSearch]
  )

  const addCondition = useCallback(
    (condition: FilterCondition<TEntity>) => {
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
    },
    [filter, setFilter]
  )

  const removeCondition = useCallback(
    (field: keyof TEntity) => {
      if (!filter) return

      const removeFromNode = (
        node: FilterNode<TEntity>
      ): FilterNode<TEntity> | undefined => {
        if (node.type === 'condition') {
          return node.field === field ? undefined : node
        } else {
          const filteredChildren = node.children
            .map(removeFromNode)
            .filter(
              (child): child is FilterNode<TEntity> => child !== undefined
            )

          if (filteredChildren.length === 0) return undefined
          if (filteredChildren.length === 1) return filteredChildren[0]

          return {
            ...node,
            children: filteredChildren
          }
        }
      }

      setFilter(removeFromNode(filter))
    },
    [filter, setFilter]
  )

  const setCondition = useCallback(
    (field: keyof TEntity, operator: FilterOperator, value: FilterValue) => {
      const condition: FilterCondition<TEntity> = {
        field,
        operator: operator,
        type: 'condition',
        value
      }

      // Remove existing condition for this field, then add new one
      removeCondition(field)
      addCondition(condition)
    },
    [addCondition, removeCondition]
  )

  const addGroup = useCallback(
    (operator: 'and' | 'or', conditions: FilterCondition<TEntity>[]) => {
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
    },
    [filter, setFilter]
  )

  const wrapInGroup = useCallback(
    (operator: 'and' | 'or') => {
      if (!filter) return

      setFilter({
        children: [filter],
        operator,
        type: 'group'
      })
    },
    [filter, setFilter]
  )

  // Pagination
  const setPage = useCallback(
    (page: number) => {
      updateSearch({
        page: {
          number: page,
          size: pageSize
        }
      })
    },
    [pageSize, updateSearch]
  )

  const setPageSize = useCallback(
    (size: number) => {
      updateSearch({
        page: {
          number: 1, // Reset to first page
          size
        }
      })
    },
    [updateSearch]
  )

  // Sorting
  const setSorting = useCallback(
    (newSorting: SortingState) => {
      updateSearch({
        sort: newSorting.map((s) => ({
          field: s.id as keyof TEntity,
          order: s.desc ? 'desc' : 'asc'
        }))
      })
      return newSorting
    },
    [updateSearch]
  )

  const addSort = useCallback(
    (field: keyof TEntity, order: 'asc' | 'desc') => {
      const newSort = { desc: order === 'desc', id: String(field) }
      setSorting([...sorting.filter((s) => s.id !== String(field)), newSort])
    },
    [sorting, setSorting]
  )

  // Reset operations
  const resetFilters = useCallback(() => {
    const { filter: _, ...rest } = searchQuery
    updateSearch(rest)
  }, [updateSearch])

  const resetPagination = useCallback(() => {
    updateSearch({ page: { number: 1, size: 10 } })
  }, [updateSearch])

  const resetSorting = useCallback(() => {
    updateSearch({ sort: [] })
  }, [updateSearch])

  const resetAll = useCallback(() => {
    setSearchQuery({})
  }, [setSearchQuery])

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
    // State
    searchQuery,
    setCondition,
    setFilter,
    setPage,
    setPageSize,
    // Actions
    setSearchQuery,
    setSorting,
    sorting,
    wrapInGroup
  }
}
