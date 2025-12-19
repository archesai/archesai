import {
  Badge,
  CalendarIcon,
  CheckIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import ExecutorForm from "#components/forms/executor-form";
import type { Executor, PageQueryParameter } from "#lib/index";
import {
  deleteExecutor,
  getListExecutorsSuspenseQueryOptions,
} from "#lib/index";

export default function ExecutorDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListExecutorsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Executor>
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
          accessorKey: "cpuShares",
          cell: ({ row }) => {
            return row.original.cpuShares;
          },
          id: "cpuShares",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "CPU Shares",
          },
        },
        {
          accessorKey: "dependencies",
          cell: ({ row }) => {
            return row.original.dependencies;
          },
          id: "dependencies",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Dependencies",
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
          accessorKey: "env",
          cell: ({ row }) => {
            return row.original.env;
          },
          id: "env",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Env",
          },
        },
        {
          accessorKey: "executeCode",
          cell: ({ row }) => {
            return row.original.executeCode;
          },
          id: "executeCode",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Execute Code",
          },
        },
        {
          accessorKey: "extraFiles",
          cell: ({ row }) => {
            return row.original.extraFiles;
          },
          id: "extraFiles",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Extra Files",
          },
        },
        {
          accessorKey: "isActive",
          cell: ({ row }) => {
            return (
              <Badge variant={row.original.isActive ? "default" : "secondary"}>
                {row.original.isActive ? "Yes" : "No"}
              </Badge>
            );
          },
          id: "isActive",
          meta: {
            filterVariant: "boolean",
            icon: CheckIcon,
            label: "Is Active",
          },
        },
        {
          accessorKey: "language",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.language}</Badge>;
          },
          id: "language",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Language",
            options: [
              { label: "Nodejs", value: "nodejs" },
              { label: "Python", value: "python" },
              { label: "Go", value: "go" },
            ],
          },
        },
        {
          accessorKey: "memoryMB",
          cell: ({ row }) => {
            return row.original.memoryMB;
          },
          id: "memoryMB",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Memory MB",
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
          accessorKey: "schemaIn",
          cell: ({ row }) => {
            return row.original.schemaIn;
          },
          id: "schemaIn",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Schema In",
          },
        },
        {
          accessorKey: "schemaOut",
          cell: ({ row }) => {
            return row.original.schemaOut;
          },
          id: "schemaOut",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Schema Out",
          },
        },
        {
          accessorKey: "timeout",
          cell: ({ row }) => {
            return row.original.timeout;
          },
          id: "timeout",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Timeout",
          },
        },
        {
          accessorKey: "version",
          cell: ({ row }) => {
            return row.original.version;
          },
          id: "version",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Version",
          },
        },
      ]}
      createForm={ExecutorForm}
      deleteItem={async (id) => {
        await deleteExecutor(id);
      }}
      entityKey="executors"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (executor) => {
        await navigate({
          params: {
            executorID: executor.id,
          },
          to: `/executors/$executorID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={ExecutorForm}
    />
  );
}
