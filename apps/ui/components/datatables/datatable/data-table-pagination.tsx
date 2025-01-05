import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { useFilterItems } from '@/hooks/useFilterItems'
import { useSelectItems } from '@/hooks/useSelectItems'
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  DoubleArrowLeftIcon,
  DoubleArrowRightIcon
} from '@radix-ui/react-icons'

interface DataTablePaginationProps<TData> {
  data: {
    metadata: {
      limit: number
      offset: number
      totalResults: number
    }
    results: TData[]
  }
}

export function DataTablePagination<TData>({
  data
}: DataTablePaginationProps<TData>) {
  const { limit, page, setLimit, setPage } = useFilterItems()
  const { selectedItems } = useSelectItems({ items: data?.results || [] })
  return (
    <div className='flex items-center justify-between'>
      {/* Display the number of items found and selected on left side*/}
      <div className='text-muted-foreground hidden text-sm sm:block'>
        {data?.metadata?.totalResults} found - {selectedItems.length || 0} of{' '}
        {Math.min(limit, data?.results.length) || 0} item(s) selected.
      </div>
      {/* Pagination controls on right side*/}
      <div className='flex items-center gap-2 lg:gap-3'>
        {/* Rows per page dropdown */}
        <div className='flex items-center gap-2'>
          <p className='hidden text-sm font-medium sm:block'>Rows per page</p>
          <Select
            onValueChange={(value) => {
              setLimit(Number(value))
            }}
            value={`${limit}`}
          >
            <SelectTrigger className='h-8 w-[70px]'>
              <SelectValue placeholder={limit} />
            </SelectTrigger>
            <SelectContent side='top'>
              {[10, 20, 30, 40, 50].map((limit) => (
                <SelectItem
                  key={limit}
                  value={`${limit}`}
                >
                  {limit}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        {/* Pagination controls */}
        <div className='flex items-center justify-center text-sm font-medium'>
          Page {page + 1} of{' '}
          {isNaN(Math.max(Math.ceil(data?.metadata?.totalResults / limit), 1))
            ? 1
            : Math.max(Math.ceil(data?.metadata?.totalResults / limit), 1)}
        </div>
        {/* Previous and Next page buttons */}
        <div className='flex items-center gap-2'>
          <Button
            className='hidden h-8 w-8 p-0 lg:flex'
            disabled={page === 0}
            onClick={() => setPage(0)}
            variant='outline'
          >
            <DoubleArrowLeftIcon className='h-5 w-5' />
          </Button>
          <Button
            className='h-8 w-8 p-0'
            disabled={page === 0}
            onClick={() => setPage(page - 1)}
            variant='outline'
          >
            <ChevronLeftIcon className='h-5 w-5' />
          </Button>
          <Button
            className='h-8 w-8 p-0'
            disabled={
              page >= Math.ceil(data?.metadata?.totalResults / limit) - 1
            }
            onClick={() => setPage(page + 1)}
            variant='outline'
          >
            <ChevronRightIcon className='h-5 w-5' />
          </Button>
          <Button
            className='hidden h-8 w-8 p-0 lg:flex'
            disabled={
              page >= Math.ceil(data?.metadata?.totalResults / limit) - 1
            }
            onClick={() =>
              setPage(Math.ceil(data?.metadata?.totalResults / limit) - 1)
            }
            variant='outline'
          >
            <DoubleArrowRightIcon className='h-5 w-5' />
          </Button>
        </div>
      </div>
    </div>
  )
}
