import type { AccessorKeyColumnDef, TableOptions } from '@tanstack/react-table'

import { useMemo, useState } from 'react'
import { getCoreRowModel, useReactTable } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/schemas'

import type { FilterActions, FilterState } from '#hooks/use-filter-state'

import { DataTableColumnHeader } from '#components/datatable/components/data-table-column-header'
import { Checkbox } from '#components/shadcn/checkbox'
import { toSentenceCase } from '#lib/utils'

interface useDataTableProps<TEntity extends BaseEntity>
  extends Omit<
    TableOptions<TEntity>,
    | 'getCoreRowModel'
    | 'manualFiltering'
    | 'manualPagination'
    | 'manualSorting'
    | 'onPaginationChange'
    | 'onSortingChange'
    | 'state'
  > {
  columns: AccessorKeyColumnDef<TEntity>[]
  filterState: FilterActions<TEntity> & FilterState<TEntity>
  total: number
}

export function useDataTable<TData extends BaseEntity>(
  props: useDataTableProps<TData>
) {
  const { columns, filterState, total, ...tableProps } = props

  // Local table state (only for UI that doesn't need URL persistence)
  const [rowSelection, setRowSelection] = useState<TData[]>([])
  // Auto-generate column headers
  const enhancedColumns = useMemo(
    () =>
      columns.map((column) => ({
        ...column,
        header:
          column.header ??
          (({ column: col }) => (
            <DataTableColumnHeader
              column={col}
              title={toSentenceCase(column.accessorKey.toString())}
            />
          ))
      })),
    [columns]
  )

  // Create table with minimal state - let filterState handle pagination/sorting
  const table = useReactTable({
    ...tableProps,
    columns: [
      // Checkbox column
      {
        cell: ({ row }) => (
          <div className='flex w-4'>
            <Checkbox
              aria-label='Select row'
              checked={rowSelection.includes(row.original)}
              onCheckedChange={(value) => {
                if (!value) {
                  setRowSelection((prev) =>
                    prev.filter((item) => item.id !== row.original.id)
                  )
                  return
                }
                setRowSelection((prev) => [...prev, row.original])
              }}
            />
          </div>
        ),
        enableHiding: false,
        enableSorting: false,
        header: () => (
          <div className='flex'>
            <Checkbox
              aria-label='Select all'
              checked={
                rowSelection === props.data ? true
                : rowSelection.length > 0 ?
                  'indeterminate'
                : false
              }
              className='translate-y-0.5'
              onCheckedChange={(value) => {
                if (value) {
                  setRowSelection((prev) => {
                    const newSelection = props.data.filter(
                      (item) => !prev.includes(item)
                    )
                    return [...prev, ...newSelection]
                  })
                } else {
                  setRowSelection([])
                }
              }}
            />
          </div>
        ),
        id: 'select'
      },
      ...enhancedColumns
    ],
    enableRowSelection: true,
    getCoreRowModel: getCoreRowModel(),
    manualFiltering: true,
    manualPagination: true,
    manualSorting: true,
    onSortingChange: (updater) => {
      const newSorting =
        typeof updater === 'function' ? updater(filterState.sorting) : updater
      filterState.setSorting(newSorting)
    },
    state: {
      pagination: {
        pageIndex: filterState.pageNumber - 1,
        pageSize: filterState.pageSize
      },
      sorting: filterState.sorting
    }
  })

  // Simple computed values
  const pageCount = Math.ceil(total / filterState.pageSize)

  return {
    pageCount,
    rows: table.getRowModel().rows,
    rowSelection,
    setRowSelection,
    sortableColumns: table
      .getAllColumns()
      .filter(
        (column) =>
          typeof column.accessorFn !== 'undefined' &&
          column.getCanSort() &&
          column.columnDef.enableSorting !== false
      ),
    table,
    totalRowCount: total
  }
}
