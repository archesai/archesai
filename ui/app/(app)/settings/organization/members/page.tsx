"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import {
  MembersControllerRemoveVariables,
  useMembersControllerFindAll,
  useMembersControllerRemove,
} from "@/generated/archesApiComponents";
import { MemberEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { endOfDay } from "date-fns";
import { User } from "lucide-react";

export default function MembersPageContent() {
  const { defaultOrgname } = useAuth();
  const { limit, page, range } = useFilterItems();

  const {
    data: members,
    isLoading,
    isPlaceholderData,
  } = useMembersControllerFindAll({
    pathParams: {
      orgname: defaultOrgname,
    },
    queryParams: {
      endDate: endOfDay(range.to || new Date()).toISOString(),
      limit,
      offset: page * limit,
      sortBy: "createdAt",
      sortDirection: "asc" as const,
      startDate: range.from?.toISOString(),
    },
  });
  const loading = isPlaceholderData || isLoading;
  const { mutateAsync: deleteMember } = useMembersControllerRemove();

  const { selectedItems } = useSelectItems({ items: members?.results || [] });

  return (
    <DataTable<MemberEntity, MembersControllerRemoveVariables>
      columns={[
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <Badge variant="outline">{row.original.role}</Badge>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Role" />
          ),
        },
        {
          accessorKey: "inviteEmail",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="font-medium">{row.original.inviteEmail}</span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Email" />
          ),
        },
        {
          accessorKey: "inviteAccepted",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="max-w-[500px] truncate font-medium">
                  {row.original.inviteAccepted}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Accepted" />
          ),
        },
      ]}
      content={() => (
        <div className="flex w-full justify-center items-center h-full">
          <User className="opacity-30" size={100} />
        </div>
      )}
      data={members as any}
      dataIcon={<User className="opacity-30" size={24} />}
      defaultView="table"
      deleteItem={deleteMember}
      getDeleteVariablesFromItem={(member) => [
        {
          pathParams: {
            memberId: member.id,
            orgname: defaultOrgname,
          },
        },
      ]}
      handleSelect={() => {}}
      itemType="API token"
      loading={loading}
      mutationVariables={selectedItems.map((id) => ({
        pathParams: {
          memberId: id,
          orgname: defaultOrgname,
        },
      }))}
    />
  );
}
