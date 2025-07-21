'use no memo'

import type { UseQueryOptions } from '@tanstack/react-query'
import type { AccessorKeyColumnDef, RowData } from '@tanstack/react-table'

import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { VisuallyHidden } from 'radix-ui'

import type { BaseEntity, SearchQuery } from '@archesai/schemas'

import type { DataTableRowAction } from '#types/simple-data-table'

import { DataTablePagination } from '#components/datatable/components/data-table-pagination'
import { TasksTableActionBar } from '#components/datatable/components/tasks-table-action-bar'
import { DataTableAdvancedToolbar } from '#components/datatable/components/toolbar/data-table-advanced-toolbar'
import { DataTableFilterMenu } from '#components/datatable/components/toolbar/data-table-filter-menu'
import { DataTableSortList } from '#components/datatable/components/toolbar/data-table-sort-list'
import { GridView } from '#components/datatable/components/views/grid-view'
import { TableView } from '#components/datatable/components/views/table-view'
// import { DataTableToolbar } from '#components/datatable/components/toolbar/data-table-toolbar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle
} from '#components/shadcn/dialog'
import { useDataTable } from '#hooks/use-data-table'
import { useFilterState } from '#hooks/use-filter-state'
import { useToggleView } from '#hooks/use-toggle-view'

declare module '@tanstack/table-core' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface TableMeta<TData extends RowData> {
    entityKey: string
    label: string
  }
}

export interface DataTableProps<TEntity extends BaseEntity> {
  actionBar?: React.ReactNode
  columns: AccessorKeyColumnDef<TEntity>[]
  createForm?: React.ComponentType
  defaultView?: 'grid' | 'table'
  deleteItem?: (id: string) => Promise<void>
  entityKey?: string
  grid?: (item: TEntity) => React.ReactNode
  gridHover?: (item: TEntity) => React.ReactNode
  handleSelect: (item: TEntity) => void
  icon: React.ReactNode
  minimal?: boolean
  updateForm?: React.ComponentType<{ id: string }>
  useFindMany: (query: SearchQuery<TEntity>) => UseQueryOptions<{
    data: TEntity[]
    meta: {
      total: number
    }
  }>
}

export function DataTable<TEntity extends BaseEntity>(
  props: DataTableProps<TEntity>
) {
  const [rowAction, setRowAction] =
    useState<DataTableRowAction<TEntity> | null>(null)

  const { setView, view } = useToggleView()
  useEffect(() => {
    setView(props.defaultView ?? 'table')
  }, [props.defaultView, setView])

  const filterState = useFilterState<TEntity>()

  const { data: queryData } = useQuery(
    props.useFindMany(filterState.searchQuery)
  )
  const data = queryData?.data ?? []
  const total = queryData?.meta.total ?? 0

  // Update table with fresh data
  const { table } = useDataTable<TEntity>({
    columns: props.columns,
    data: data,
    filterState,
    pageCount: Math.ceil(total / (filterState.searchQuery.page?.size ?? 10)) // Should come from backend
  })

  return (
    <div className='flex flex-1 flex-col gap-4'>
      {/* FILTER TOOLBAR */}
      <DataTableAdvancedToolbar table={table}>
        <DataTableSortList
          align='start'
          table={table}
        />
        <DataTableFilterMenu table={table} />
      </DataTableAdvancedToolbar>
      {/* {!props.minimal && <DataTableToolbar<TEntity> table={table} />} */}

      {/* DATA TABLE - EITHER GRID OR TABLE VIEW*/}
      <div className='flex-1 overflow-auto'>
        {view === 'grid' ?
          <GridView<TEntity>
            icon={props.icon}
            table={table}
          />
        : <TableView<TEntity> table={table} />}
      </div>

      {/* PAGINATION */}
      {!props.minimal && <DataTablePagination<TEntity> table={table} />}

      {/* THIS IS THE FORM DIALOG */}
      <Dialog
        onOpenChange={() => {
          setRowAction(null)
        }}
        open={
          rowAction?.variant === 'update' || rowAction?.variant === 'custom'
        }
      >
        <VisuallyHidden.Root>
          <DialogDescription />
          <DialogTitle>
            {rowAction?.variant === 'update' ? 'Edit' : 'Create'}{' '}
            {table.options.meta?.entityKey ?? 'Entity'}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DialogContent
          aria-description='Create/Edit'
          className='p-0'
          title='Create/Edit'
        >
          {rowAction?.variant === 'update' && props.updateForm && (
            <props.updateForm id={rowAction.row.original.id} />
          )}

          {rowAction?.variant === 'create' && props.createForm && (
            <props.createForm />
          )}
        </DialogContent>
      </Dialog>
      {<TasksTableActionBar table={table} />}
    </div>
  )
}
