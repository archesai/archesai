"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { RunStatusButton } from "@/components/run-status-button";
import {
  RunsControllerFindAllPathParams,
  useRunsControllerFindAll,
} from "@/generated/archesApiComponents";
import { PipelineRunEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { PackageCheck } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function RunsPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <DataTable<PipelineRunEntity, RunsControllerFindAllPathParams, undefined>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="max-w-[200px] shrink truncate font-medium text-primary"
                  href={`/runs/single?runId=${row.original.id}`}
                >
                  {row.original.name}
                </Link>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
        {
          accessorKey: "status",
          cell: ({ row }) => {
            return <RunStatusButton run={row.original} />;
          },
          enableHiding: false,
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Status"
            />
          ),
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <span className="font-light">
                {format(
                  new Date(row.original.createdAt),
                  "yyyy-MM-dd HH:mm:ss"
                )}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Created"
            />
          ),
        },
      ]}
      content={() => (
        <div className="flex h-full w-full items-center justify-center"></div>
      )}
      dataIcon={<PackageCheck size={24} />}
      defaultView="table"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={() => {}}
      handleSelect={(run) => router.push(`/tools/runs/single?run=${run.id}`)}
      itemType="run"
      useFindAll={useRunsControllerFindAll}
      useRemove={() => {
        return { mutateAsync: async () => {} };
      }}
    />
  );
}
