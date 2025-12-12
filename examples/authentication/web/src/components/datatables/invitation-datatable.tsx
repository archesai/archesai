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
import InvitationForm from "#components/forms/invitation-form";
import type {
  Invitation,
  InvitationsFilterParameter,
  InvitationsSortParameter,
  PageQueryParameter,
} from "#lib/index";
import {
  deleteInvitation,
  getListInvitationsSuspenseQueryOptions,
} from "#lib/index";

export default function InvitationDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListInvitationsSuspenseQueryOptions({
      filter: query.filter as unknown as InvitationsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as InvitationsSortParameter,
    });
  };

  return (
    <DataTableContainer<Invitation>
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
          accessorKey: "email",
          cell: ({ row }) => {
            return row.original.email;
          },
          id: "email",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Email",
          },
        },
        {
          accessorKey: "expiresAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.expiresAt} />;
          },
          id: "expiresAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Expires At",
          },
        },
        {
          accessorKey: "inviterID",
          cell: ({ row }) => {
            return row.original.inviterID;
          },
          id: "inviterID",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Inviter ID",
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
          accessorKey: "status",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.status}</Badge>;
          },
          id: "status",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Status",
            options: [
              { label: "Pending", value: "pending" },
              { label: "Accepted", value: "accepted" },
              { label: "Declined", value: "declined" },
              { label: "Expired", value: "expired" },
            ],
          },
        },
      ]}
      createForm={InvitationForm}
      deleteItem={async (id) => {
        await deleteInvitation(id);
      }}
      entityKey="invitations"
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (invitation) => {
        await navigate({
          params: {
            invitationID: invitation.id,
          },
          to: `/invitations/$invitationID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={InvitationForm}
    />
  );
}
