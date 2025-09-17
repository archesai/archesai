import type {
  PageQueryParameter,
  Run,
  RunsFilterParameter,
  RunsSortParameter,
} from "@archesai/client";
import { deleteRun, getListRunsSuspenseQueryOptions } from "@archesai/client";
import {
  PackageCheckIcon,
  StatusTypeEnumButton,
  Timestamp,
} from "@archesai/ui";
import { RUN_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import { DataTableContainer } from "#components/datatables/data-table-container";
import { Link } from "#components/navigation/link";

export default function RunDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListRunsSuspenseQueryOptions({
      filter: query.filter as unknown as RunsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as RunsSortParameter,
    });
  };

  return (
    <DataTableContainer<Run>
      columns={[
        {
          accessorKey: "id",
          cell: ({ row }) => {
            return (
              <Link
                className="max-w-[200px] shrink truncate font-medium"
                params={{
                  runID: row.original.id,
                }}
                to={`/runs/$runID`}
              >
                {row.original.id}
              </Link>
            );
          },
          id: "id",
        },
        {
          accessorKey: "status",
          cell: ({ row }) => {
            return (
              <StatusTypeEnumButton
                run={row.original}
                size="sm"
              />
            );
          },
          id: "status",
        },
        {
          accessorKey: "duration",
          cell: ({ row }) => {
            return row.original.startedAt && row.original.completedAt ? (
              <Timestamp
                date={new Date(
                  new Date(row.original.completedAt).getTime() -
                    new Date(row.original.startedAt).getTime(),
                ).toISOString()}
              />
            ) : (
              "N/A"
            );
          },
          id: "duration",
        },
        {
          accessorKey: "startedAt",
          cell: ({ row }) => {
            return row.original.startedAt ? (
              <Timestamp date={row.original.startedAt} />
            ) : (
              "N/A"
            );
          },
          id: "startedAt",
        },
        {
          accessorKey: "completedAt",
          cell: ({ row }) => {
            return row.original.completedAt ? (
              <Timestamp date={row.original.completedAt} />
            ) : (
              "N/A"
            );
          },
          id: "completedAt",
        },
      ]}
      deleteItem={async (id) => {
        await deleteRun(id);
      }}
      entityKey={RUN_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (run) => {
        await navigate({
          params: {
            runID: run.id,
          },
          to: `/runs/$runID`,
        });
      }}
      icon={<PackageCheckIcon />}
    />
  );
}
