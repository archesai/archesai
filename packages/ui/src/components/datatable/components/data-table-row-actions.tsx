'use no memo'

import type { Table } from '@tanstack/table-core'

import { MoreHorizontal } from 'lucide-react'

import type { BaseEntity } from '@archesai/schemas'

import { DeleteItems } from '#components/custom/delete-items'
import { Button } from '#components/shadcn/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '#components/shadcn/dropdown-menu'

export interface DataTableRowActionsProps<TEntity extends BaseEntity> {
  deleteItem?: (id: string) => Promise<void>
  getEditFormFromItem?: (item: TEntity) => React.ReactNode
  row: { original: TEntity }
  setFinalForm?: (form: React.ReactNode) => void
  setFormOpen?: (open: boolean) => void
  table: Table<TEntity>
}

export function DataTableRowActions<TEntity extends BaseEntity>(
  props: DataTableRowActionsProps<TEntity>
) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label='Expand row options'
          className='flex h-8 w-8 p-0 data-[state=open]:bg-muted'
          variant='ghost'
        >
          <MoreHorizontal className='h-5 w-5' />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align='end'
        className='w-[160px]'
      >
        {props.getEditFormFromItem ?
          <>
            <DropdownMenuItem
              onClick={() => {
                if (!props.getEditFormFromItem) {
                  throw new Error('getEditFormFromItem function is not defined')
                }
                if (props.setFinalForm && props.setFormOpen) {
                  props.setFinalForm(
                    props.getEditFormFromItem(props.row.original)
                  )
                  props.setFormOpen(true)
                }
              }}
            >
              Edit
            </DropdownMenuItem>
            <DropdownMenuSeparator />
          </>
        : null}
        {props.deleteItem ?
          <DropdownMenuItem
            onSelect={(e) => {
              e.preventDefault()
            }}
          >
            <DeleteItems
              deleteItem={async (id) => {
                if (!props.deleteItem) {
                  throw new Error('deleteItem function is not defined')
                }
                await props.deleteItem(id)
                props.table.toggleAllRowsSelected(false)
              }}
              entityKey={props.table.options.meta?.entityKey ?? 'Entity'}
              items={[props.row.original]}
              variant='md'
            />
          </DropdownMenuItem>
        : null}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
