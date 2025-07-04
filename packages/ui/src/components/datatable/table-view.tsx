import type { Table as ReactTable } from '@tanstack/react-table'

import { flexRender } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/domain'

import type { DataTableContainerProps } from '#components/datatable/data-table'

import { Skeleton } from '#components/shadcn/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '#components/shadcn/table'

export interface TableViewProps<TEntity extends BaseEntity>
  extends Pick<
    DataTableContainerProps<TEntity>,
    | 'columns'
    | 'data'
    | 'deleteItem'
    | 'entityType'
    | 'isFetched'
    | 'selectedItems'
  > {
  table: ReactTable<TEntity>
}

export function TableView<TEntity extends BaseEntity>({
  columns,
  entityType,
  isFetched,
  table
}: TableViewProps<TEntity>) {
  return (
    <div className='rounded-lg border shadow-xs'>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder ? null : (
                    flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )
                  )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {!isFetched ?
            <>
              {Array.from({ length: 10 }).map((_, index) => (
                <TableRow key={index}>
                  {columns.map((_column, i) => (
                    <TableCell
                      className='h-12 p-2'
                      key={i}
                    >
                      <Skeleton className='h-4' />
                    </TableCell>
                  ))}
                  <TableCell className='h-12 p-2'>
                    <Skeleton className='h-4' />
                  </TableCell>
                </TableRow>
              ))}
            </>
          : table.getRowModel().rows.length ?
            table.getRowModel().rows.map((row) => (
              <TableRow
                data-state={row.getIsSelected() && 'selected'}
                key={row.id}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          : <TableRow>
              <TableCell
                className='h-24 text-center'
                colSpan={columns.length + 2}
              >
                No {entityType} found
              </TableCell>
            </TableRow>
          }
        </TableBody>
      </Table>
    </div>
  )
}
