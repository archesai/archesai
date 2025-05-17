'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { format } from 'date-fns'
import { PackageCheck, Workflow } from 'lucide-react'

import type { RunEntity } from '@archesai/domain'

import { deleteRun, useFindManyRuns } from '@archesai/client'
import { RUN_ENTITY_KEY } from '@archesai/domain'
import { StatusTypeEnumButton } from '@archesai/ui/components/custom/run-status-button'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Button } from '@archesai/ui/components/shadcn/button'

export default function RunDataTable() {
  const router = useRouter()

  return (
    <DataTable<RunEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <Link
                  className='max-w-[200px] shrink truncate font-medium'
                  href={`/run/single?runId=${row.original.id}`}
                >
                  {row.original.name}
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
                {row.original.runType === 'PIPELINE_RUN' ? (
                  <Workflow className='h-5 w-5' />
                ) : (
                  <PackageCheck className='h-5 w-5' />
                )}
              </Button>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'status',
          cell: ({ row }) => {
            return (
              <div className='pl-3'>
                <StatusTypeEnumButton
                  run={row.original}
                  size='sm'
                />
              </div>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'duration',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.startedAt && row.original.completedAt
                  ? format(
                      new Date(row.original.completedAt).getTime() -
                        new Date(row.original.startedAt).getTime(),
                      'mm:ss'
                    )
                  : 'N/A'}
              </span>
            )
          },
          enableHiding: false,
          enableSorting: false
        },
        {
          accessorKey: 'startedAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.startedAt
                  ? format(new Date(row.original.startedAt), 'M/d/yy h:mm a')
                  : 'N/A'}
              </span>
            )
          }
        },
        {
          accessorKey: 'completedAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.completedAt
                  ? format(new Date(row.original.completedAt), 'M/d/yy h:mm a')
                  : 'N/A'}
              </span>
            )
          }
        }
      ]}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteRun(id)
      }}
      entityType={RUN_ENTITY_KEY}
      handleSelect={(run) => {
        router.push(`/run/single?runId=${run.id}`)
      }}
      icon={<PackageCheck />}
      useFindMany={useFindManyRuns}
    />
  )
}
