"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import {
  PipelinesControllerFindAllPathParams,
  PipelinesControllerRemoveVariables,
  usePipelinesControllerFindAll,
  usePipelinesControllerRemove,
} from "@/generated/archesApiComponents";
import { PipelineEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { Workflow } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ContentPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <DataTable<
      PipelineEntity,
      PipelinesControllerFindAllPathParams,
      PipelinesControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex gap-2">
                <Link
                  className="max-w-[200px] shrink truncate font-medium text-primary"
                  href={`/pipelines/single?pipelineId=${row.original.id}`}
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
          accessorKey: "description",
          cell: ({ row }) => {
            return (
              <span>
                {(row.original.description || "No Description").toString()}
              </span>
            );
          },
          enableHiding: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="text-sm"
              column={column}
              title="Description"
            />
          ),
        },
        {
          accessorKey: "transformations",
          cell: ({ row }) => {
            return (
              <div className="flex gap-1">
                {row.original.pipelineSteps?.map((step) => {
                  return (
                    <Badge className="text-primary" variant="outline">
                      {step.tool}
                    </Badge>
                  );
                })}
              </div>
            );
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
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <span className="font-light">
                {format(new Date(row.original.createdAt), "M/d/yy h:mm a")}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
          ),
        },
      ]}
      dataIcon={<Workflow />}
      defaultView="table"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(pipeline) => ({
        pathParams: {
          orgname: defaultOrgname,
          pipelineId: pipeline.id,
        },
      })}
      handleSelect={(pipeline) =>
        router.push(`/pipelines/single?pipelineId=${pipeline.id}`)
      }
      itemType="pipeline"
      useFindAll={usePipelinesControllerFindAll}
      useRemove={usePipelinesControllerRemove}
    />
  );
}
