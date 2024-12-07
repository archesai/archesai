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

interface TableViewProps<TItem extends BaseItem> {
  columns: any[] // Replace with appropriate type if possible
  items?: TItem[]
  itemType: string
  table: any // Replace with appropriate type from react-table if possible
}

export function TableView<TItem extends BaseItem>({
  columns,
  itemType,
  table
}: TableViewProps<TItem>) {
  return (
    <div className='rounded-md border bg-sidebar shadow-sm'>
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
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row: any, index: number) => (
              <TableRow
                className={
                  'transition-all hover:bg-muted' +
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
