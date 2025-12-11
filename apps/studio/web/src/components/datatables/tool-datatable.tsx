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
import ToolForm from "#components/forms/tool-form";
import type { PageQueryParameter, Tool } from "#lib/index";
import { deleteTool, getListToolsSuspenseQueryOptions } from "#lib/index";

export default function ToolDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListToolsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Tool>
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
          accessorKey: "inputMimeType",
          cell: ({ row }) => {
            return row.original.inputMimeType;
          },
          id: "inputMimeType",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Input Mime Type",
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
        {
          accessorKey: "outputMimeType",
          cell: ({ row }) => {
            return row.original.outputMimeType;
          },
          id: "outputMimeType",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Output Mime Type",
          },
        },
      ]}
      createForm={ToolForm}
      deleteItem={async (id) => {
        await deleteTool(id);
      }}
      entityKey="tools"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (tool) => {
        await navigate({
          params: {
            toolID: tool.id,
          },
          to: `/tools/$toolID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={ToolForm}
    />
  );
}
