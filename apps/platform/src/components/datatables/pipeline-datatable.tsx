import { Link, useNavigate } from '@tanstack/react-router'

import type { PipelineEntity } from '@archesai/domain'

import { deletePipeline, useFindManyPipelines } from '@archesai/client'
import { PIPELINE_ENTITY_KEY } from '@archesai/domain'
import { Workflow } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

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
                  className='max-w-[200px] shrink truncate font-medium text-primary'
                  params={{ pipelineId: row.original.id }}
                  to={`/pipelines/$pipelineId`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          }
        },
        {
          accessorKey: 'description',
          cell: ({ row }) => {
            return (
              <span>
                {(row.original.description || 'No Description').toString()}
              </span>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'Inputs',
          cell: ({ row }) => {
            return (
              <div className='flex gap-1'>
                {row.original.steps.map((step, i) => {
                  return <Badge key={i}>{step.tool.name}</Badge>
                })}
              </div>
            )
          },
          enableHiding: false,
          enableSorting: false
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          }
        }
      ]}
      defaultView='table'
      deleteItem={async (id) => {
        await deletePipeline(id)
      }}
      entityType={PIPELINE_ENTITY_KEY}
      handleSelect={async (pipeline) => {
        await navigate({ to: `/pipelines/single?pipelineId=${pipeline.id}` })
      }}
      icon={<Workflow />}
      useFindMany={useFindManyPipelines}
    />
  )
}
