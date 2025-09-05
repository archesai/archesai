'use no memo'

import type { Table } from '@tanstack/react-table'
import type { JSX } from 'react'

import type { BaseEntity } from '#types/entities'

import {
  CheckIcon,
  ChevronsUpDownIcon,
  Settings2Icon
} from '#components/custom/icons'
import { Button } from '#components/shadcn/button'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '#components/shadcn/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '#components/shadcn/popover'
import { cn, toSentenceCase } from '#lib/utils'

interface DataTableViewOptionsProps<TEntity extends BaseEntity> {
  table: Table<TEntity>
}

export function DataTableViewOptions<TEntity extends BaseEntity>(
  props: DataTableViewOptionsProps<TEntity>
): JSX.Element {
  const columns = props.table
    .getAllColumns()
    .filter(
      (column) =>
        typeof column.accessorFn !== 'undefined' && column.getCanHide()
    )

  return (
    <Popover modal>
      <PopoverTrigger asChild>
        <Button
          aria-label='Toggle columns'
          className='ml-auto hidden h-8 lg:flex'
          role='combobox'
          size='sm'
          variant='ghost'
        >
          <Settings2Icon />
          View
          <ChevronsUpDownIcon className='ml-auto opacity-50' />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        align='end'
        className='w-44 p-0'
      >
        <Command>
          <CommandInput placeholder='Search columns...' />
          <CommandList>
            <CommandEmpty>No columns found.</CommandEmpty>
            <CommandGroup>
              {columns.map((column) => {
                return (
                  <CommandItem
                    key={column.id}
                    onSelect={() => {
                      column.toggleVisibility(!column.getIsVisible())
                    }}
                  >
                    <span className='truncate'>
                      {column.columnDef.meta?.label ??
                        toSentenceCase(column.id)}
                    </span>
                    <CheckIcon
                      className={cn(
                        'ml-auto size-4 shrink-0',
                        column.getIsVisible() ? 'opacity-100' : 'opacity-0'
                      )}
                    />
                  </CommandItem>
                )
              })}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
