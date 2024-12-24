// components/datatable/TableView.tsx
'use client'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table'
import { flexRender } from '@tanstack/react-table'

import { BaseItem } from './data-table'
import { Skeleton } from '@/components/ui/skeleton'

interface TableViewProps<TItem extends BaseItem> {
  columns: any[] // Replace with appropriate type if possible
  items?: TItem[]
  itemType: string
  table: any // Replace with appropriate type from react-table if possible
  isFetched: boolean
}

export function TableView<TItem extends BaseItem>({
  columns,
  itemType,
  table,
  isFetched
}: TableViewProps<TItem>) {
  return (
    <div className='rounded-lg border shadow-sm'>
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup: any) => (
            <TableRow
              className='bg-background/40'
              key={headerGroup.id}
            >
              {headerGroup.headers.map((header: any, i: number) => (
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
              {[...Array(10)].map((_, index) => (
                <TableRow
                  className={
                    'hover:bg-muted transition-all' +
                    (index % 2 ? ' bg-background/40' : ' ')
                  }
                  key={index}
                >
                  {columns.map((column, i) => (
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
          ) : table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row: any, index: number) => (
              <TableRow
                className={
                  'hover:bg-muted transition-all' +
                  (index % 2 ? ' bg-background/40' : ' ')
                }
                data-state={row.getIsSelected() && 'selected'}
                key={row.id}
              >
                {row.getVisibleCells().map((cell: any) => (
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
                No {itemType}s found
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
