import {
  Badge,
  CalendarIcon,
  PackageCheckIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import { TOOL_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import { DataTableContainer } from "#components/datatables/data-table-container";
import { Link } from "#components/navigation/link";
import type {
  PageQueryParameter,
  Tool,
  ToolsFilterParameter,
  ToolsSortParameter,
} from "#lib/index";
import { deleteTool, getListToolsSuspenseQueryOptions } from "#lib/index";

export default function ToolDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListToolsSuspenseQueryOptions({
      filter: query.filter as unknown as ToolsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as ToolsSortParameter,
    });
  };

  return (
    <DataTableContainer<Tool>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <Link
                className="shrink truncate text-wrap text-blue-500 hover:underline md:text-sm"
                params={{
                  artifactID: row.original.id,
                }}
                to={`/artifacts/$artifactID`}
              >
                {row.original.name}
              </Link>
            );
          },
          enableColumnFilter: true,
          id: "name",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Name",
          },
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return row.original.description || "No Description";
          },
          enableColumnFilter: true,
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
            return (
              <Badge variant={"secondary"}>{row.original.inputMimeType}</Badge>
            );
          },
          enableColumnFilter: true,
          id: "inputMimeType",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Input",
          },
        },
        {
          accessorKey: "outputMimeType",
          cell: ({ row }) => {
            return (
              <Badge variant={"secondary"}>{row.original.outputMimeType}</Badge>
            );
          },
          enableColumnFilter: true,
          id: "outputMimeType",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Output",
          },
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          enableColumnFilter: true,
          id: "createdAt",
          meta: {
            filterVariant: "dateRange",
            icon: CalendarIcon,
            label: "Created at",
          },
        },
      ]}
      deleteItem={async (id) => {
        await deleteTool(id);
      }}
      entityKey={TOOL_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (tool) => {
        await navigate({
          to: `/tool/single?toolID=${tool.id}`,
        });
      }}
      icon={<PackageCheckIcon />}
    />
  );
}
