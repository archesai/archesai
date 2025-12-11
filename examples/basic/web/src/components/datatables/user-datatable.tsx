import {
  Badge,
  CalendarIcon,
  CheckIcon,
  DataTableContainer,
  ListIcon,
  TextIcon,
  Timestamp,
} from "@archesai/ui";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import UserForm from "#components/forms/user-form";
import type { PageQueryParameter, User } from "#lib/index";
import { deleteUser, getListUsersSuspenseQueryOptions } from "#lib/index";

export default function UserDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListUsersSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<User>
      columns={[
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            const val = row.original.createdAt;
            return val ? <Timestamp date={val as string} /> : "-";
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
            const val = row.original.updatedAt;
            return val ? <Timestamp date={val as string} /> : "-";
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
          accessorKey: "emailVerified",
          cell: ({ row }) => {
            return (
              <Badge
                variant={row.original.emailVerified ? "default" : "secondary"}
              >
                {row.original.emailVerified ? "Yes" : "No"}
              </Badge>
            );
          },
          id: "emailVerified",
          meta: {
            filterVariant: "boolean",
            icon: CheckIcon,
            label: "Email Verified",
          },
        },
        {
          accessorKey: "image",
          cell: ({ row }) => {
            return row.original.image;
          },
          id: "image",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Image",
          },
        },
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.name}</Badge>;
          },
          id: "name",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Name",
          },
        },
      ]}
      deleteItem={async (id) => {
        await deleteUser(id);
      }}
      entityKey="users"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (user) => {
        await navigate({
          params: {
            userID: user.id,
          },
          to: `/users/$userID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={UserForm}
    />
  );
}
