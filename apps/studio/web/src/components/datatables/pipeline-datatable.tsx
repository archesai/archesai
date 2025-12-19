import {
  Badge,
  CalendarIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import PipelineForm from "#components/forms/pipeline-form";
import type { Pipeline } from "#lib/index";
import {
  deletePipeline,
  getListPipelinesSuspenseQueryOptions,
} from "#lib/index";

export default function PipelineDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (_query: SearchQuery) => {
    return getListPipelinesSuspenseQueryOptions();
  };

  return (
    <DataTableContainer<Pipeline>
      columns={[
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            const val = row.original.createdAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "createdAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Created At",
          },
        },
        {
          accessorKey: "updatedAt",
          cell: ({ row }) => {
            const val = row.original.updatedAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "updatedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Updated At",
          },
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return row.original.description;
          },
          id: "description",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Description",
          },
        },
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.name}</Badge>;
          },
          id: "name",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Name",
          },
        },
        {
          accessorKey: "organizationID",
          cell: ({ row }) => {
            return row.original.organizationID;
          },
          id: "organizationID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Organization ID",
          },
        },
      ]}
      createForm={PipelineForm}
      deleteItem={async (id) => {
        await deletePipeline(id);
      }}
      entityKey="pipelines"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (pipeline) => {
        await navigate({
          params: {
            pipelineID: pipeline.id,
          },
          to: `/pipelines/$pipelineID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={PipelineForm}
    />
  );
}
