import type {
  PageQueryParameter,
  Pipeline,
  PipelinesFilterParameter,
  PipelinesSortParameter,
} from "@archesai/client";
import {
  deletePipeline,
  getListPipelinesSuspenseQueryOptions,
} from "@archesai/client";
import { WorkflowIcon } from "@archesai/ui/components/custom/icons";
import { Timestamp } from "@archesai/ui/components/custom/timestamp";
import { DataTable } from "@archesai/ui/components/datatable/data-table";
import { PIPELINE_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { Link, useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";

export default function PipelineDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListPipelinesSuspenseQueryOptions({
      filter: query.filter as unknown as PipelinesFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as PipelinesSortParameter,
    });
  };

  return (
    <DataTable<Pipeline>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="max-w-[200px] shrink truncate font-medium text-blue-500"
                  params={{ pipelineID: row.original.id }}
                  to={`/pipelines/$pipelineID`}
                >
                  {row.original.name}
                </Link>
              </div>
            );
          },
          id: "name",
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return row.original.description ?? "No Description";
          },
          id: "description",
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          id: "createdAt",
        },
      ]}
      deleteItem={async (id) => {
        await deletePipeline(id);
      }}
      entityKey={PIPELINE_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (pipeline) => {
        await navigate({
          params: {
            pipelineID: pipeline.id,
          },
          to: `/pipelines/$pipelineID`,
        });
      }}
      icon={<WorkflowIcon />}
    />
  );
}
