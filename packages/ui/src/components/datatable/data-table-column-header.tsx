import type { Column } from '@tanstack/react-table'

import { ArrowDown, ArrowsUpFromLine, ArrowUp } from 'lucide-react'

import { Button } from '#components/shadcn/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '#components/shadcn/dropdown-menu'
import { cn } from '#lib/utils'

interface DataTableColumnHeaderProps<TData, TValue>
  extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>
  title: string
}

export function DataTableColumnHeader<TData, TValue>({
  className,
  column,
  title
}: DataTableColumnHeaderProps<TData, TValue>) {
  if (!column.getCanSort()) {
    return <div className={cn(className, 'text-xs')}>{title}</div>
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          className='-ml-3 h-8 data-[state=open]:bg-muted'
          size='sm'
          variant='ghost'
        >
          <span>{title}</span>
          {column.getIsSorted() === 'desc' ?
            <ArrowDown className='h-4 w-4' />
          : column.getIsSorted() === 'asc' ?
            <ArrowUp className='h-4 w-4' />
          : <ArrowsUpFromLine className='h-4 w-4' />}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align='start'>
        <DropdownMenuItem
          className='flex items-center gap-2'
          onClick={() => {
            column.toggleSorting(false)
          }}
        >
          <ArrowUp className='h-4 w-4 text-muted-foreground/70' />
          <span>Asc</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          className='flex items-center gap-2'
          onClick={() => {
            column.toggleSorting(true)
          }}
        >
          <ArrowDown className='h-4 w-4 text-muted-foreground/70' />
          <span>Desc</span>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          className='flex items-center gap-2'
          onClick={() => {
            column.toggleVisibility(false)
          }}
        >
          <ArrowsUpFromLine className='h-4 w-4 text-muted-foreground/70' />
          <span>Hide</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
