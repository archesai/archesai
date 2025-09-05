'use no memo'

import type { Table } from '@tanstack/react-table'
import type { JSX } from 'react'

import * as React from 'react'

// import { toast } from 'sonner'

import type { BaseEntity } from '#types/entities'

import { DownloadIcon, TrashIcon } from '#components/custom/icons'
import {
  DataTableActionBar,
  DataTableActionBarAction,
  DataTableActionBarSelection
} from '#components/datatable/components/data-table-action-bar'
import { Separator } from '#components/shadcn/separator'
import { exportTableToCSV } from '#lib/export'

const _actions = ['export', 'delete'] as const

type Action = (typeof _actions)[number]

interface TableActionBarProps<TEntity extends BaseEntity> {
  table: Table<TEntity>
}

export function TasksTableActionBar<TEntity extends BaseEntity>({
  table
}: TableActionBarProps<TEntity>): JSX.Element {
  const rows = table.getFilteredSelectedRowModel().rows
  const [isPending, startTransition] = React.useTransition()
  const [currentAction, setCurrentAction] = React.useState<Action | null>(null)

  const getIsActionPending = React.useCallback(
    (action: Action) => isPending && currentAction === action,
    [isPending, currentAction]
  )

  // const onUpdate = React.useCallback(
  //   ({
  //     field,
  //     value
  //   }: {
  //     field: 'priority' | 'status'
  //     value: Task['priority'] | Task['status']
  //   }) => {
  //     setCurrentAction(field === 'status' ? 'update-status' : 'update-priority')
  //     startTransition(async () => {
  //       const { error } = await updateTasks({
  //         [field]: value,
  //         ids: rows.map((row) => row.original.id)
  //       })

  //       if (error) {
  //         toast.error(error)
  //         return
  //       }
  //       toast.success('Tasks updated')
  //     })
  //   },
  //   [rows]
  // )

  const onExport = React.useCallback(() => {
    setCurrentAction('export')
    startTransition(() => {
      exportTableToCSV(table, {
        excludeColumns: ['select', 'actions'],
        onlySelected: true
      })
    })
  }, [table])

  // const onDelete = React.useCallback(() => {
  //   setCurrentAction('delete')
  //   startTransition(async () => {
  //     const { error } = await deleteTasks({
  //       ids: rows.map((row) => row.original.id)
  //     })

  //     if (error) {
  //       toast.error(error)
  //       return
  //     }
  //     table.toggleAllRowsSelected(false)
  //   })
  // }, [rows, table])

  return (
    <DataTableActionBar
      table={table}
      visible={rows.length > 0}
    >
      <DataTableActionBarSelection table={table} />
      <Separator
        className='hidden data-[orientation=vertical]:h-5 sm:block'
        orientation='vertical'
      />
      <div className='flex items-center gap-1.5'>
        {/* <Select
          onValueChange={(value: Task['status']) => {
            onTaskUpdate({ field: 'status', value })
          }}
        >
          <SelectTrigger asChild>
            <DataTableActionBarAction
              isPending={getIsActionPending('update-status')}
              size='icon'
              tooltip='Update status'
            >
              <CheckCircle2 />
            </DataTableActionBarAction>
          </SelectTrigger>
          <SelectContent align='center'>
            <SelectGroup>
              {tasks.status.enumValues.map((status) => (
                <SelectItem
                  className='capitalize'
                  key={status}
                  value={status}
                >
                  {status}
                </SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
        <Select
          onValueChange={(value: Task['priority']) => {
            onTaskUpdate({ field: 'priority', value })
          }}
        >
          <SelectTrigger asChild>
            <DataTableActionBarAction
              isPending={getIsActionPending('update-priority')}
              size='icon'
              tooltip='Update priority'
            >
              <ArrowUp />
            </DataTableActionBarAction>
          </SelectTrigger>
          <SelectContent align='center'>
            <SelectGroup>
              {tasks.priority.enumValues.map((priority) => (
                <SelectItem
                  className='capitalize'
                  key={priority}
                  value={priority}
                >
                  {priority}
                </SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select> */}
        <DataTableActionBarAction
          isPending={getIsActionPending('export')}
          onClick={onExport}
          size='icon'
          tooltip='Export'
        >
          <DownloadIcon />
        </DataTableActionBarAction>
        <DataTableActionBarAction
          isPending={getIsActionPending('delete')}
          // onClick={onDelete}
          size='icon'
          tooltip='Delete'
        >
          <TrashIcon />
        </DataTableActionBarAction>
      </div>
    </DataTableActionBar>
  )
}
