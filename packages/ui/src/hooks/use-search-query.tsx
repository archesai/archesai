import { useAtom } from 'jotai'

import { searchQueryAtom } from '#atoms/search-query'

type SortDirection = 'asc' | 'desc'

export const useSearchQuery = () => {
  const [searchQuery, setSearchQuery] = useAtom(searchQueryAtom)

  // Sorting
  const sortField = searchQuery?.sort?.replace(/^-/, '') ?? 'createdAt'
  const sortDirection: SortDirection = searchQuery?.sort?.startsWith('-')
    ? 'desc'
    : 'asc'

  const updateSort = (field: string, direction: SortDirection) => {
    setSearchQuery((prev) => ({
      ...prev,
      sort: `${direction === 'desc' ? '-' : ''}${field}`
    }))
  }

  const setSortBy = (field: string) => {
    updateSort(field, sortDirection)
  }
  const setSortDirection = (direction: SortDirection) => {
    updateSort(sortField, direction)
  }

  // Pagination
  const pageNumber = searchQuery?.page?.number ?? 1
  const pageSize = searchQuery?.page?.size ?? 10

  const setPage = (number: number, size?: number) => {
    setSearchQuery((prev) => ({
      ...prev,
      page: {
        number,
        size: size ?? prev?.page?.size ?? 10
      }
    }))
  }

  // Filtering
  const filter = searchQuery?.filter ?? {}

  const setFilter = (filter: object) => {
    setSearchQuery((prev) => ({
      ...prev,
      filter: {
        ...prev?.filter,
        ...filter
      }
    }))
  }

  const resetQuery = () => {
    setSearchQuery({
      filter: {},
      page: { number: 1, size: 10 },
      sort: 'createdAt'
    })
  }

  return {
    // Filter
    filter,
    // Pagination
    page: searchQuery?.page,
    pageNumber,

    pageSize,
    resetQuery,
    searchQuery,
    setFilter,
    setPage,

    setSearchQuery,
    setSortBy,
    setSortDirection,
    // Sort
    sort: searchQuery?.sort ?? 'createdAt',

    sortBy: sortField,
    sortDirection
  }
}
