import { Link, useNavigate } from '@tanstack/react-router'

import type { PipelineEntity } from '@archesai/schemas'

import {
  deletePipeline,
  getFindManyPipelinesSuspenseQueryOptions
} from '@archesai/client'
import { PIPELINE_ENTITY_KEY } from '@archesai/schemas'
import { WorkflowIcon } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'

export default function PipelineDataTable() {
  const navigate = useNavigate()

  return (
    <DataTable<PipelineEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <Link
                  className='max-w-[200px] shrink truncate font-medium text-blue-500'
                  params={{ pipelineId: row.original.id }}
                  to={`/pipelines/$pipelineId`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          },
          id: 'name'
        },
        {
          accessorKey: 'description',
          cell: ({ row }) => {
            return row.original.description ?? 'No Description'
          },
          id: 'description'
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: 'createdAt'
        }
      ]}
      deleteItem={async (id) => {
        await deletePipeline(id)
      }}
      entityKey={PIPELINE_ENTITY_KEY}
      getQueryOptions={getFindManyPipelinesSuspenseQueryOptions}
      handleSelect={async (pipeline) => {
        await navigate({
          params: {
            pipelineId: pipeline.id
          },
          to: `/pipelines/$pipelineId`
        })
      }}
      icon={<WorkflowIcon />}
    />
  )
}
