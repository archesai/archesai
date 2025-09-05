import type { JSX } from 'react'

import { Link, useNavigate } from '@tanstack/react-router'

import type { FindManyRunsParams } from '@archesai/client'
import type { RunEntity, SearchQuery } from '@archesai/ui/types/entities'

import {
  deleteRun,
  getFindManyRunsSuspenseQueryOptions
} from '@archesai/client'
import { PackageCheckIcon } from '@archesai/ui/components/custom/icons'
import { StatusTypeEnumButton } from '@archesai/ui/components/custom/run-status-button'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { RUN_ENTITY_KEY } from '@archesai/ui/lib/constants'

export default function RunDataTable(): JSX.Element {
  const navigate = useNavigate()

  const getQueryOptions = (query: SearchQuery) => {
    const params: any =
      query.filter || query.page || query.sort ?
        {
          ...(query.filter && {
            filter: query.filter as unknown as FindManyRunsParams['filter']
          }),
          ...(query.page && { page: query.page }),
          ...(query.sort && { sort: query.sort as FindManyRunsParams['sort'] })
        }
      : undefined
    return getFindManyRunsSuspenseQueryOptions(params) as any
  }

  return (
    <DataTable<RunEntity>
      columns={[
        {
          accessorKey: 'id',
          cell: ({ row }) => {
            return (
              <Link
                className='max-w-[200px] shrink truncate font-medium'
                params={{
                  runId: row.original.id
                }}
                to={`/runs/$runId`}
              >
                {row.original.id}
              </Link>
            )
          },
          id: 'id'
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
          id: 'status'
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
          id: 'duration'
        },
        {
          accessorKey: 'startedAt',
          cell: ({ row }) => {
            return row.original.startedAt ?
                <Timestamp date={row.original.startedAt} />
              : 'N/A'
          },
          id: 'startedAt'
        },
        {
          accessorKey: 'completedAt',
          cell: ({ row }) => {
            return row.original.completedAt ?
                <Timestamp date={row.original.completedAt} />
              : 'N/A'
          },
          id: 'completedAt'
        }
      ]}
      deleteItem={async (id) => {
        await deleteRun(id)
      }}
      entityKey={RUN_ENTITY_KEY}
      getQueryOptions={getQueryOptions}
      handleSelect={async (run) => {
        await navigate({ params: { runId: run.id }, to: `/runs/$runId` })
      }}
      icon={<PackageCheckIcon />}
    />
  )
}
