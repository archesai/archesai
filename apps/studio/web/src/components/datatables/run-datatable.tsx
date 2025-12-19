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
import RunForm from "#components/forms/run-form";
import type { PageQueryParameter, Run } from "#lib/index";
import { deleteRun, getListRunsSuspenseQueryOptions } from "#lib/index";

export default function RunDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListRunsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Run>
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
          accessorKey: "completedAt",
          cell: ({ row }) => {
            const val = row.original.completedAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "completedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Completed At",
          },
        },
        {
          accessorKey: "error",
          cell: ({ row }) => {
            return row.original.error;
          },
          id: "error",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Error",
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
        {
          accessorKey: "pipelineID",
          cell: ({ row }) => {
            return row.original.pipelineID;
          },
          id: "pipelineID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Pipeline ID",
          },
        },
        {
          accessorKey: "progress",
          cell: ({ row }) => {
            return row.original.progress;
          },
          id: "progress",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Progress",
          },
        },
        {
          accessorKey: "startedAt",
          cell: ({ row }) => {
            const val = row.original.startedAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "startedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Started At",
          },
        },
        {
          accessorKey: "status",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.status}</Badge>;
          },
          id: "status",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Status",
            options: [
              { label: "COMPLETED", value: "COMPLETED" },
              { label: "FAILED", value: "FAILED" },
              { label: "PROCESSING", value: "PROCESSING" },
              { label: "QUEUED", value: "QUEUED" },
            ],
          },
        },
        {
          accessorKey: "toolID",
          cell: ({ row }) => {
            return row.original.toolID;
          },
          id: "toolID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Tool ID",
          },
        },
      ]}
      createForm={RunForm}
      deleteItem={async (id) => {
        await deleteRun(id);
      }}
      entityKey="runs"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (run) => {
        await navigate({
          params: {
            runID: run.id,
          },
          to: `/runs/$runID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={RunForm}
    />
  );
}
