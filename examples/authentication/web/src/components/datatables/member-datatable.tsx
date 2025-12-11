import {
  Badge,
  CalendarIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import MemberForm from "#components/forms/member-form";
import type {
  Member,
  MembersFilterParameter,
  MembersSortParameter,
  PageQueryParameter,
} from "#lib/index";
import { deleteMember, getListMembersSuspenseQueryOptions } from "#lib/index";

export default function MemberDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListMembersSuspenseQueryOptions({
      filter: query.filter as unknown as MembersFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as MembersSortParameter,
    });
  };

  return (
    <DataTableContainer<Member>
      columns={[
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          id: "createdAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Created At",
          },
        },
        {
          accessorKey: "updatedAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.updatedAt} />;
          },
          id: "updatedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Updated At",
          },
        },
        {
          accessorKey: "organizationID",
          cell: ({ row }) => {
            return row.original.organizationID;
          },
          id: "organizationID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Organization ID",
          },
        },
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.role}</Badge>;
          },
          id: "role",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Role",
            options: [
              { label: "Admin", value: "admin" },
              { label: "Owner", value: "owner" },
              { label: "Basic", value: "basic" },
            ],
          },
        },
        {
          accessorKey: "userID",
          cell: ({ row }) => {
            return row.original.userID;
          },
          id: "userID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "User ID",
          },
        },
      ]}
      createForm={MemberForm}
      deleteItem={async (id) => {
        await deleteMember(id);
      }}
      entityKey="members"
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (member) => {
        await navigate({
          params: {
            memberID: member.id,
          },
          to: `/members/$memberID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={MemberForm}
    />
  );
}
