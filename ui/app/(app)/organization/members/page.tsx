"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import MemberForm from "@/components/forms/member-form";
import { Badge } from "@/components/ui/badge";
import {
  MembersControllerFindAllPathParams,
  MembersControllerRemoveVariables,
  useMembersControllerFindAll,
  useMembersControllerRemove,
} from "@/generated/archesApiComponents";
import { MemberEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/use-auth";
import { CheckIcon, User, XIcon } from "lucide-react";

export default function MembersPageContent() {
  const { defaultOrgname } = useAuth();

  return (
    <DataTable<
      MemberEntity,
      MembersControllerFindAllPathParams,
      MembersControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return (
              <Badge className="text-primary" variant="outline">
                {row.original.role}
              </Badge>
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
              <span className="font-medium">{row.original.inviteEmail}</span>
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
              <span className="max-w-[500px] truncate font-medium">
                {row.original.inviteAccepted ? (
                  <CheckIcon className="text-primary" />
                ) : (
                  <XIcon className="text-red-950" />
                )}
              </span>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Accepted" />
          ),
        },
      ]}
      content={() => (
        <div className="flex h-full w-full items-center justify-center">
          <User className="opacity-30" size={100} />
        </div>
      )}
      createForm={<MemberForm />}
      dataIcon={<User className="opacity-30" size={24} />}
      defaultView="table"
      filterField="inviteEmail"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(member) => ({
        pathParams: {
          memberId: member.id,
          orgname: defaultOrgname,
        },
      })}
      getEditFormFromItem={(member) => {
        return <MemberForm memberId={member.id} />;
      }}
      handleSelect={() => {}}
      itemType="Member"
      useFindAll={useMembersControllerFindAll}
      useRemove={useMembersControllerRemove}
    />
  );
}
