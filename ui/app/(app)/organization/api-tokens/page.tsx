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
              <div className="flex space-x-2">
                <Badge
                  variant={
                    row.original.role === "ADMIN" ? "default" : "secondary"
                  }
                >
                  {row.original.role}
                </Badge>
              </div>
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
              <div className="flex space-x-2">
                <span className="max-w-[500px] truncate font-medium">
                  {row.original.name}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
        {
          accessorKey: "key",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="font-medium">{row.original.key}</span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Value" />
          ),
        },
      ]}
      content={() => (
        <div className="flex w-full justify-center items-center h-full">
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
