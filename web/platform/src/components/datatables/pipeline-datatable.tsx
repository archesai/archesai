import type { JSX } from "react"

import { Link, useNavigate } from "@tanstack/react-router"

import type {
  PageQueryParameter,
  PipelineEntity,
  PipelinesFilterParameter,
  PipelinesSortParameter
} from "@archesai/client"
import type { SearchQuery } from "@archesai/ui/types/entities"

import {
  deletePipeline,
  getFindManyPipelinesSuspenseQueryOptions
} from "@archesai/client"
import { WorkflowIcon } from "@archesai/ui/components/custom/icons"
import { Timestamp } from "@archesai/ui/components/custom/timestamp"
import { DataTable } from "@archesai/ui/components/datatable/data-table"
import { PIPELINE_ENTITY_KEY } from "@archesai/ui/lib/constants"

export default function PipelineDataTable(): JSX.Element {
  const navigate = useNavigate()

  const getQueryOptions = (query: SearchQuery) => {
    return getFindManyPipelinesSuspenseQueryOptions({
      filter: query.filter as unknown as PipelinesFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as PipelinesSortParameter
    })
  }

  return (
    <DataTable<PipelineEntity>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="max-w-[200px] shrink truncate font-medium text-blue-500"
                  params={{ pipelineId: row.original.id }}
                  to={`/pipelines/$pipelineId`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          },
          id: "name"
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return row.original.description ?? "No Description"
          },
          id: "description"
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: "createdAt"
        }
      ]}
      deleteItem={async (id) => {
        await deletePipeline(id)
      }}
      entityKey={PIPELINE_ENTITY_KEY}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
      getQueryOptions={getQueryOptions as any}
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
