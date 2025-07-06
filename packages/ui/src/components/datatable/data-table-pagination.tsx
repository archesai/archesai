import type { Table } from '@tanstack/react-table'

import { ArrowBigLeftDash, ArrowBigRightDash } from 'lucide-react'

import type { BaseEntity } from '@archesai/domain'

import { Button } from '#components/shadcn/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '#components/shadcn/select'
import { useSearchQuery } from '#hooks/use-search-query'
import { cn } from '#lib/utils'

interface DataTablePaginationProps<TEntity extends BaseEntity> {
  response: {
    data: TEntity[]
    meta: {
      page: number
      size: number
      total_records: number
    }
  }
  table: Table<TEntity>
}

export function DataTablePagination<TEntity extends BaseEntity>({
  response,
  table
}: DataTablePaginationProps<TEntity>) {
  const { pageNumber, pageSize, searchQuery, setPage, setSearchQuery } =
    useSearchQuery()
  const selected = table.getSelectedRowModel().rows
  return (
    <div
      className={cn(
        'flex w-full flex-col-reverse items-center justify-between gap-4 overflow-auto p-1 sm:flex-row sm:gap-8'
      )}
    >
      {/* Display the number of items found and selected on left side*/}
      <div className='flex-1 text-sm whitespace-nowrap text-muted-foreground'>
        {response.meta.total_records} found - {selected.length} of{' '}
        {Math.min(pageSize, response.data.length)} item(s) selected.
      </div>
      {/* Pagination controls on right side*/}
      <div className='flex flex-col-reverse items-center gap-4 sm:flex-row sm:gap-6 lg:gap-8'>
        {/* Rows per page dropdown */}
        <div className='flex items-center space-x-2'>
          <p className='text-sm font-medium whitespace-nowrap'>Rows per page</p>
          <Select
            onValueChange={(value) => {
              setSearchQuery({
                ...searchQuery,
                filter: {
                  ...searchQuery?.filter
                },
                page: {
                  ...searchQuery?.page,
                  size: parseInt(value)
                }
              })
            }}
            value={pageSize.toString()}
          >
            <SelectTrigger
              className='h-8 w-[4.5rem] [&[data-size]]:h-8'
              size='sm'
            >
              <SelectValue placeholder={pageSize.toString()} />
            </SelectTrigger>
            <SelectContent side='top'>
              {[10, 20, 30, 40, 50].map((limit) => (
                <SelectItem
                  key={limit}
                  value={limit.toString()}
                >
                  {limit}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        {/* Pagination controls */}
        <div className='flex items-center justify-center text-sm font-medium'>
          Page {pageNumber + 1} of{' '}
          {(
            isNaN(
              Math.max(Math.ceil(response.meta.total_records / pageSize), 1)
            )
          ) ?
            1
          : Math.max(Math.ceil(response.meta.total_records / pageSize), 1)}
        </div>
        {/* Previous and Next page buttons */}
        <div className='flex items-center gap-2'>
          <Button
            className='hidden lg:flex'
            disabled={pageSize === 0}
            onClick={() => {
              setPage(0, pageSize)
            }}
            size={'sm'}
            variant='outline'
          >
            <ArrowBigLeftDash className='h-5 w-5' />
          </Button>
          <Button
            disabled={pageSize === 0}
            onClick={() => {
              setPage(pageNumber - 1, pageSize)
            }}
            size='sm'
            variant='outline'
          >
            <ArrowBigLeftDash className='h-5 w-5' />
          </Button>
          <Button
            disabled={
              pageSize >= Math.ceil(response.meta.total_records / pageSize) - 1
            }
            onClick={() => {
              setPage(pageNumber + 1, pageSize)
            }}
            size='sm'
            variant='outline'
          >
            <ArrowBigRightDash className='h-5 w-5' />
          </Button>
          <Button
            className='hidden lg:flex'
            disabled={
              pageSize >= Math.ceil(response.meta.total_records / pageSize) - 1
            }
            onClick={() => {
              setPage(
                Math.ceil(response.meta.total_records / pageSize),
                pageSize
              )
            }}
            size='sm'
            variant='outline'
          >
            <ArrowBigRightDash className='h-5 w-5' />
          </Button>
        </div>
      </div>
    </div>
  )
}
