'use no memo'

import type { Table } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/schemas'

import { DataTableViewOptions } from '#components/datatable/components/data-table-view-options'
import { ViewToggle } from '#components/datatable/components/view-toggle'
import { cn } from '#lib/utils'

interface DataTableAdvancedToolbarProps<TEntity extends BaseEntity>
  extends React.ComponentProps<'div'> {
  table: Table<TEntity>
}

export function DataTableAdvancedToolbar<TEntity extends BaseEntity>({
  children,
  className,
  table,
  ...props
}: DataTableAdvancedToolbarProps<TEntity>) {
  return (
    <div
      aria-orientation='horizontal'
      className={cn(
        'flex w-full items-start justify-between gap-2 p-1',
        className
      )}
      role='toolbar'
      {...props}
    >
      <div className='flex flex-1 flex-wrap items-center gap-2'>{children}</div>
      <div className='flex items-center gap-2'>
        <ViewToggle />
        <DataTableViewOptions table={table} />
      </div>
    </div>
  )
}
