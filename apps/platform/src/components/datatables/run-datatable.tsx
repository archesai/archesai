import { Link, useNavigate } from '@tanstack/react-router'

import type { RunEntity } from '@archesai/schemas'

import {
  deleteRun,
  getFindManyRunsSuspenseQueryOptions
} from '@archesai/client'
import { RUN_ENTITY_KEY } from '@archesai/schemas'
import { PackageCheck, Workflow } from '@archesai/ui/components/custom/icons'
import { StatusTypeEnumButton } from '@archesai/ui/components/custom/run-status-button'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Button } from '@archesai/ui/components/shadcn/button'

export default function RunDataTable() {
  const navigate = useNavigate()

  return (
    <DataTable<RunEntity>
      columns={[
        {
          accessorKey: 'id',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <Link
                  className='max-w-[200px] shrink truncate font-medium'
                  params={{
                    runId: row.original.id
                  }}
                  to={`/runs/$runId`}
                >
                  {row.original.id}
                </Link>
              </div>
            )
          }
        },
        {
          accessorKey: 'runType',
          cell: ({ row }) => {
            return (
              <Button
                className='text-primary hover:text-primary/90'
                size='sm'
                variant={'outline'}
              >
                {row.original.runType === 'PIPELINE_RUN' ?
                  <Workflow className='h-5 w-5' />
                : <PackageCheck className='h-5 w-5' />}
              </Button>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'status',
          cell: ({ row }) => {
            return (
              <StatusTypeEnumButton
                run={row.original}
                size='sm'
              />
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'duration',
          cell: ({ row }) => {
            return row.original.startedAt && row.original.completedAt ?
                <Timestamp
                  date={new Date(
                    new Date(row.original.completedAt).getTime() -
                      new Date(row.original.startedAt).getTime()
                  ).toISOString()}
                />
              : 'N/A'
          },
          enableHiding: false,
          enableSorting: false
        },
        {
          accessorKey: 'startedAt',
          cell: ({ row }) => {
            return row.original.startedAt ?
                <Timestamp date={row.original.startedAt} />
              : 'N/A'
          }
        },
        {
          accessorKey: 'completedAt',
          cell: ({ row }) => {
            return row.original.completedAt ?
                <Timestamp date={row.original.completedAt} />
              : 'N/A'
          }
        }
      ]}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteRun(id)
      }}
      entityKey={RUN_ENTITY_KEY}
      handleSelect={async (run) => {
        await navigate({ params: { runId: run.id }, to: `/runs/$runId` })
      }}
      icon={<PackageCheck />}
      useFindMany={getFindManyRunsSuspenseQueryOptions}
    />
  )
}
