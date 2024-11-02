"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { JobStatusButton } from "@/components/job-status-button";
import {
  JobsControllerFindAllPathParams,
  JobsControllerRemoveVariables,
  useJobsControllerFindAll,
  useJobsControllerRemove,
} from "@/generated/archesApiComponents";
import { JobEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { PackageCheck } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ToolsJobsPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <DataTable<
      JobEntity,
      JobsControllerFindAllPathParams,
      JobsControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="max-w-[200px] shrink truncate font-medium text-primary"
                  href={`/tool/single?toolId=${row.original.id}`}
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
            return <JobStatusButton job={row.original} />;
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
          accessorKey: "input",
          cell: ({ row }) => {
            return <span>{row.original.input}</span>;
          },
          enableHiding: false,
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Input"
            />
          ),
        },
        {
          accessorKey: "output",
          cell: ({ row }) => {
            return <span>{row.original.input}</span>;
          },
          enableHiding: false,
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Output"
            />
          ),
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <span className="font-light">
                {new Date(row.original.createdAt).toLocaleDateString()}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
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
      getDeleteVariablesFromItem={(job) => ({
        pathParams: {
          id: job.id,
          orgname: defaultOrgname,
        },
      })}
      handleSelect={(job) => router.push(`/job/single?jobId=${job.id}`)}
      itemType="job"
      useFindAll={useJobsControllerFindAll}
      useRemove={useJobsControllerRemove}
    />
  );
}
