"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import {
  ToolsControllerFindAllPathParams,
  ToolsControllerRemoveVariables,
  useToolsControllerFindAll,
  useToolsControllerRemove,
} from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { PackageCheck } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ContentPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <DataTable<
      ToolEntity,
      ToolsControllerFindAllPathParams,
      ToolsControllerRemoveVariables
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
          accessorKey: "description",
          cell: ({ row }) => {
            return <span>{row.original.description || "No Description"}</span>;
          },
          enableHiding: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Description"
            />
          ),
        },
        {
          accessorKey: "inputType",
          cell: ({ row }) => {
            return (
              <Badge className="text-primary" variant={"secondary"}>
                {row.original.inputType}
              </Badge>
            );
          },
          enableHiding: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="-ml-2 text-sm"
              column={column}
              title="Input"
            />
          ),
        },
        {
          accessorKey: "outputType",
          cell: ({ row }) => {
            return (
              <Badge className="text-primary" variant={"secondary"}>
                {row.original.outputType}
              </Badge>
            );
          },
          enableHiding: false,
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
                {format(new Date(row.original.createdAt), "M/d/yy h:mm a")}
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
      dataIcon={<PackageCheck />}
      defaultView="table"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(tool) => ({
        pathParams: {
          orgname: defaultOrgname,
          toolId: tool.id,
        },
      })}
      handleSelect={(tool) => router.push(`/tool/single?toolId=${tool.id}`)}
      itemType="tool"
      useFindAll={useToolsControllerFindAll}
      useRemove={useToolsControllerRemove}
    />
  );
}
