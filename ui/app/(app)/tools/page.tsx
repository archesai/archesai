"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import {
  ToolsControllerFindAllPathParams,
  ToolsControllerRemoveVariables,
  useToolsControllerFindAll,
  useToolsControllerRemove,
} from "@/generated/archesApiComponents";
import { ToolEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { File } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ContentPage() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();

  return (
    <div className="flex h-full flex-col gap-3">
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
              return (
                <span>{row.original.description || "No Description"}</span>
              );
            },
            enableHiding: false,
            enableSorting: false,
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
              return <span>{row.original.inputType}</span>;
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
            accessorKey: "outputType",
            cell: ({ row }) => {
              return <span>{row.original.outputType}</span>;
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
        dataIcon={<File size={24} />}
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
    </div>
  );
}
