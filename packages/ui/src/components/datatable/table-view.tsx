'use client'

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
    <div className='shadow-xs rounded-lg border'>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow
              className='bg-background/40'
              key={headerGroup.id}
            >
              {headerGroup.headers.map((header, i: number) => (
                <TableHead
                  className={'text-base' + (i === 0 ? ' w-4' : '')}
                  colSpan={header.colSpan}
                  key={header.id}
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {!isFetched ? (
            <>
              {Array.from({ length: 10 }).map((_, index) => (
                <TableRow
                  className={
                    'hover:bg-muted transition-all' +
                    (index % 2 ? ' bg-background/40' : ' ')
                  }
                  key={index}
                >
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
          ) : table.getRowModel().rows.length ? (
            table.getRowModel().rows.map((row, index: number) => (
              <TableRow
                className={
                  'hover:bg-muted transition-all' +
                  (index % 2 ? ' bg-background/40' : ' ')
                }
                data-state={row.getIsSelected() && 'selected'}
                key={row.id}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell
                    className='p-2'
                    key={cell.id}
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell
                className='h-24 text-center'
                colSpan={columns.length + 2}
              >
                No {entityType}s found
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
