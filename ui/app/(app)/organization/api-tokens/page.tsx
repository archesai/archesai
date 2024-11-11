"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import APITokenForm from "@/components/forms/api-token-form";
import { Badge } from "@/components/ui/badge";
import {
  ApiTokensControllerFindAllPathParams,
  ApiTokensControllerRemoveVariables,
  useApiTokensControllerFindAll,
  useApiTokensControllerRemove,
} from "@/generated/archesApiComponents";
import { ApiTokenEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { format } from "date-fns";
import { User } from "lucide-react";
import { useRouter } from "next/navigation";

export default function ApiTokensPageContent() {
  const { defaultOrgname } = useAuth();
  const router = useRouter();

  return (
    <DataTable<
      ApiTokenEntity,
      ApiTokensControllerFindAllPathParams,
      ApiTokensControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return (
              <Badge className="text-primary" variant="secondary">
                {row.original.role}
              </Badge>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Role" />
          ),
        },
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <span className="max-w-[500px] truncate font-medium">
                {row.original.name}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
        {
          accessorKey: "key",
          cell: ({ row }) => {
            return <span className="font-medium">{row.original.key}</span>;
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Value" />
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
      content={() => (
        <div className="flex h-full w-full items-center justify-center">
          <User className="opacity-30" size={100} />
        </div>
      )}
      createForm={<APITokenForm />}
      dataIcon={<User className="opacity-30" size={24} />}
      defaultView="table"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(apiToken) => ({
        pathParams: {
          id: apiToken.id,
          orgname: defaultOrgname,
        },
      })}
      getEditFormFromItem={(apiToken) => (
        <APITokenForm apiTokenId={apiToken.id} />
      )}
      handleSelect={(apiToken) => router.push(`/apiTokens/${apiToken.id}/chat`)}
      itemType="API token"
      useFindAll={useApiTokensControllerFindAll}
      useRemove={useApiTokensControllerRemove}
    />
  );
}
