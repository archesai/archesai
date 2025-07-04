import type { Table } from '@tanstack/react-table'

import { CrossIcon, GridIcon, ListIcon } from 'lucide-react'

import type { BaseEntity } from '@archesai/domain'

import type { DataTableContainerProps } from '#components/datatable/data-table'

import { DataTableViewOptions } from '#components/datatable/data-table-view-options'
import { DatePickerWithRange } from '#components/datatable/date-range-picker'
import { Button } from '#components/shadcn/button'
import { Checkbox } from '#components/shadcn/checkbox'
import { useSelectItems } from '#hooks/use-select-items'
import { useToggleView } from '#hooks/use-toggle-view'

export interface DataTableToolbarProps<TEntity extends BaseEntity>
  extends Pick<
    DataTableContainerProps<TEntity>,
    'createForm' | 'data' | 'entityType' | 'readonly' | 'setFormOpen'
  > {
  table: Table<TEntity>
}

export function DataTableToolbar<TEntity extends BaseEntity>({
  createForm,
  data,
  entityType,
  readonly,
  setFormOpen,
  table
}: DataTableToolbarProps<TEntity>) {
  const isFiltered = table.getState().columnFilters.length > 0

  const { selectedAllItems, selectedSomeItems, toggleSelectAll } =
    useSelectItems({ items: data })

  // const { searchQuery, setSearchQuery } = useSearchQuery()
  return (
    <div className='flex flex-wrap items-center gap-2'>
      {!readonly && (
        <Checkbox
          aria-label='Select all'
          checked={selectedAllItems || (selectedSomeItems && 'indeterminate')}
          onCheckedChange={() => {
            toggleSelectAll()
          }}
        />
      )}

      {isFiltered && (
        <Button
          onClick={() => {
            table.resetColumnFilters()
          }}
          size='sm'
          variant='outline'
        >
          <span>Reset</span>
          <CrossIcon className='h-5 w-5' />
        </Button>
      )}

      <DatePickerWithRange />
      <ViewToggle />
      <DataTableViewOptions table={table} />
      {createForm && !readonly ?
        <Button
          onClick={() => {
            setFormOpen(true)
          }}
          size='sm'
          variant={'outline'}
        >
          Create {entityType.toLowerCase()}
        </Button>
      : null}
    </div>
  )
}

export function ViewToggle() {
  const { setView, view } = useToggleView()
  return (
    <div className='hidden h-8 gap-2 md:flex'>
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === 'table' ?
            'bg-secondary text-primary'
          : 'bg-transparent text-muted-foreground'
        }`}
        onClick={() => {
          setView('table')
        }}
        size='icon'
        variant={'outline'}
      >
        <ListIcon className='h-5 w-5' />
      </Button>
      <Button
        className={`flex h-full items-center justify-center transition-colors ${
          view === 'grid' ?
            'bg-secondary text-primary'
          : 'bg-transparent text-muted-foreground'
        }`}
        onClick={() => {
          setView('grid')
        }}
        size='icon'
        variant={'outline'}
      >
        <GridIcon className='h-5 w-5' />
      </Button>
    </div>
  )
}
