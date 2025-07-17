'use no memo'

import type { UseQueryOptions } from '@tanstack/react-query'
import type { AccessorKeyColumnDef, RowData } from '@tanstack/react-table'

import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { VisuallyHidden } from 'radix-ui'

import type { BaseEntity, SearchQuery } from '@archesai/schemas'

import { DataTablePagination } from '#components/datatable/components/data-table-pagination'
import { GridView } from '#components/datatable/components/grid-view'
import { TableView } from '#components/datatable/components/table-view'
import { TasksTableActionBar } from '#components/datatable/components/tasks-table-action-bar'
import { DataTableAdvancedToolbar } from '#components/datatable/components/toolbar/data-table-advanced-toolbar'
import { DataTableFilterMenu } from '#components/datatable/components/toolbar/data-table-filter-menu'
import { DataTableSortList } from '#components/datatable/components/toolbar/data-table-sort-list'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle
} from '#components/shadcn/dialog'
import { useDataTable } from '#hooks/use-data-table'
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
  createForm?: React.ReactNode
  defaultView?: 'grid' | 'table'
  deleteItem?: (id: string) => Promise<void>
  entityKey?: string
  getEditFormFromItem?: (item: TEntity) => React.ReactNode
  grid?: (item: TEntity) => React.ReactNode
  gridHover?: (item: TEntity) => React.ReactNode
  handleSelect: (item: TEntity) => void
  icon: React.ReactNode
  minimal?: boolean
  readonly?: boolean
  useFindMany: (query: SearchQuery<TEntity>) => UseQueryOptions<{
    data: TEntity[]
  }>
}

export function DataTable<TEntity extends BaseEntity>(
  props: DataTableProps<TEntity>
) {
  // Use the useDebounce hook to debounce the query
  // const debouncedQuery = useDebounce(query, 200) // 500ms delay

  const [formOpen, setFormOpen] = useState(false)
  const [finalForm, setFinalForm] = useState<React.ReactNode | undefined>(
    props.createForm
  )

  const { setView, view } = useToggleView()
  useEffect(() => {
    setView(props.defaultView ?? 'table')
  }, [props.defaultView, setView])

  const dataTableResult = useDataTable<TEntity>({
    columns: props.columns,
    data: [], // Will be filled from query
    pageCount: -1 // Placeholder, should come from backend response
  })

  const { searchQuery } = dataTableResult

  const { data: queryData } = useQuery(props.useFindMany(searchQuery))
  const data = queryData?.data ?? []

  // Update table with fresh data
  const { table } = useDataTable<TEntity>({
    columns: props.columns,
    data: data,
    pageCount: Math.ceil(1000 / (searchQuery.page?.size ?? 10)) // Should come from backend
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
      {!props.minimal && (
        <div className='self-auto'>
          <DataTablePagination<TEntity> table={table} />
        </div>
      )}

      {/* THIS IS THE FORM DIALOG */}
      <Dialog
        onOpenChange={(o) => {
          setFormOpen(o)
          if (!o) {
            setFinalForm(props.createForm)
          }
        }}
        open={formOpen}
      >
        <VisuallyHidden.Root>
          <DialogDescription />
          <DialogTitle>
            {finalForm ? 'Edit' : 'Create'}{' '}
            {table.options.meta?.entityKey ?? 'Entity'}
          </DialogTitle>
        </VisuallyHidden.Root>
        <DialogContent
          aria-description='Create/Edit'
          className='p-0'
          title='Create/Edit'
        >
          {finalForm}
        </DialogContent>
      </Dialog>
      {table.getFilteredSelectedRowModel().rows.length > 0 && (
        <TasksTableActionBar table={table} />
      )}
    </div>
  )
}
