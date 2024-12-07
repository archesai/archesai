'use client'
import { DataTable } from '@/components/datatables/datatable/data-table'
import { DataTableColumnHeader } from '@/components/datatables/datatable/data-table-column-header'
import { RunStatusButton } from '@/components/run-status-button'
import {
  RunsControllerFindAllPathParams,
  RunsControllerRemoveVariables,
  useRunsControllerFindAll,
  useRunsControllerRemove
} from '@/generated/archesApiComponents'
import { RunEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { format } from 'date-fns'
import { PackageCheck, Workflow } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

import { Button } from '../ui/button'

export default function RunDataTable() {
  const router = useRouter()
  const { defaultOrgname } = useAuth()

  return (
    <DataTable<
      RunEntity & {
        name: string
      },
      RunsControllerFindAllPathParams,
      RunsControllerRemoveVariables
    >
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
          },
          header: ({ column }) => <DataTableColumnHeader column={column} title='Name' />
        },
        {
          accessorKey: 'runType',
          cell: ({ row }) => {
            return (
              <Button className='text-primary hover:text-primary/90' size='sm' variant={'outline'}>
                {row.original.runType === 'PIPELINE_RUN' ? (
                  <Workflow className='h-5 w-5' />
                ) : (
                  <PackageCheck className='h-5 w-5' />
                )}
              </Button>
            )
          },
          enableHiding: false,
          header: ({ column }) => <DataTableColumnHeader className='-ml-2 text-sm' column={column} title='Run Type' />
        },
        {
          accessorKey: 'status',
          cell: ({ row }) => {
            return (
              <div className='pl-3'>
                <RunStatusButton run={row.original} size='sm' />
              </div>
            )
          },
          enableHiding: false,
          header: ({ column }) => <DataTableColumnHeader className='-ml-2 text-sm' column={column} title='Input' />
        },
        {
          accessorKey: 'duration',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.startedAt && row.original.completedAt
                  ? format(
                      new Date(row.original.completedAt).getTime() - new Date(row.original.startedAt).getTime(),
                      'mm:ss'
                    )
                  : 'N/A'}
              </span>
            )
          },
          enableHiding: false,
          enableSorting: false,
          header: ({ column }) => <DataTableColumnHeader className='-ml-2 text-sm' column={column} title='Duration' />
        },
        {
          accessorKey: 'startedAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.startedAt ? format(new Date(row.original.startedAt), 'M/d/yy h:mm a') : 'N/A'}
              </span>
            )
          },
          header: ({ column }) => <DataTableColumnHeader className='-ml-2 text-sm' column={column} title='Started At' />
        },
        {
          accessorKey: 'completedAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {row.original.completedAt ? format(new Date(row.original.completedAt), 'M/d/yy h:mm a') : 'N/A'}
              </span>
            )
          },
          header: ({ column }) => (
            <DataTableColumnHeader className='-ml-2 text-sm' column={column} title='Completed At' />
          )
        }
      ]}
      dataIcon={<PackageCheck />}
      defaultView='table'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(run) => ({
        pathParams: {
          id: run.id,
          orgname: defaultOrgname
        }
      })}
      handleSelect={(run) => router.push(`/run/single?runId=${run.id}`)}
      itemType='run'
      useFindAll={useRunsControllerFindAll}
      useRemove={useRunsControllerRemove}
    />
  )
}
