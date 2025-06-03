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
import { useSelectItems } from '#hooks/use-select-items'

interface DataTablePaginationProps<TEntity extends BaseEntity> {
  response: {
    data: TEntity[]
    meta: {
      page: number
      size: number
      total_records: number
    }
  }
}

export function DataTablePagination<TEntity extends BaseEntity>({
  response
}: DataTablePaginationProps<TEntity>) {
  const { pageNumber, pageSize, searchQuery, setPage, setSearchQuery } =
    useSearchQuery()
  const { selectedItems } = useSelectItems({ items: [] })
  return (
    <div className='flex items-center justify-between'>
      {/* Display the number of items found and selected on left side*/}
      <div className='hidden text-sm text-muted-foreground sm:block'>
        {response.meta.total_records} found - {selectedItems.length || 0} of{' '}
        {Math.min(pageSize, response.data.length) || 0} item(s) selected.
      </div>
      {/* Pagination controls on right side*/}
      <div className='flex items-center gap-2 lg:gap-3'>
        {/* Rows per page dropdown */}
        <div className='flex items-center gap-2'>
          <p className='hidden text-sm font-medium sm:block'>Rows per page</p>
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
            <SelectTrigger className='h-8 w-[70px]'>
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
            className='hidden h-8 w-8 p-0 lg:flex'
            disabled={pageSize === 0}
            onClick={() => {
              setPage(0, pageSize)
            }}
            variant='outline'
          >
            <ArrowBigLeftDash className='h-5 w-5' />
          </Button>
          <Button
            className='h-8 w-8 p-0'
            disabled={pageSize === 0}
            onClick={() => {
              setPage(pageNumber - 1, pageSize)
            }}
            variant='outline'
          >
            <ArrowBigLeftDash className='h-5 w-5' />
          </Button>
          <Button
            className='h-8 w-8 p-0'
            disabled={
              pageSize >= Math.ceil(response.meta.total_records / pageSize) - 1
            }
            onClick={() => {
              setPage(pageNumber + 1, pageSize)
            }}
            variant='outline'
          >
            <ArrowBigRightDash className='h-5 w-5' />
          </Button>
          <Button
            className='hidden h-8 w-8 p-0 lg:flex'
            disabled={
              pageSize >= Math.ceil(response.meta.total_records / pageSize) - 1
            }
            onClick={() => {
              setPage(
                Math.ceil(response.meta.total_records / pageSize),
                pageSize
              )
            }}
            variant='outline'
          >
            <ArrowBigRightDash className='h-5 w-5' />
          </Button>
        </div>
      </div>
    </div>
  )
}
