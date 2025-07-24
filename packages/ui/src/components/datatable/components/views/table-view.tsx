'use no memo'

import type { Table as ReactTable } from '@tanstack/react-table'

import { flexRender } from '@tanstack/react-table'

import type { BaseEntity } from '@archesai/schemas'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '#components/shadcn/table'

export interface TableViewProps<TEntity extends BaseEntity> {
  table: ReactTable<TEntity>
}

export function TableView<TEntity extends BaseEntity>(
  props: TableViewProps<TEntity>
) {
  const columns = props.table.getAllColumns()
  return (
    <div className='border bg-card'>
      <Table>
        <TableHeader>
          {props.table.getHeaderGroups().map((headerGroup) => (
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
          {props.table.getRowModel().rows.length ?
            props.table.getRowModel().rows.map((row) => (
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
                No items found
              </TableCell>
            </TableRow>
          }
        </TableBody>
      </Table>
    </div>
  )
}
