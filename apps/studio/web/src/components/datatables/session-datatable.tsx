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
import SessionForm from "#components/forms/session-form";
import type { PageQueryParameter, Session } from "#lib/index";
import { deleteSession, getListSessionsSuspenseQueryOptions } from "#lib/index";

export default function SessionDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListSessionsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Session>
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
          accessorKey: "authMethod",
          cell: ({ row }) => {
            return row.original.authMethod;
          },
          id: "authMethod",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Auth Method",
          },
        },
        {
          accessorKey: "authProvider",
          cell: ({ row }) => {
            return (
              <Badge variant="secondary">{row.original.authProvider}</Badge>
            );
          },
          id: "authProvider",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Auth Provider",
            options: [
              { label: "Local", value: "local" },
              { label: "Google", value: "google" },
              { label: "Github", value: "github" },
              { label: "Microsoft", value: "microsoft" },
              { label: "Apple", value: "apple" },
            ],
          },
        },
        {
          accessorKey: "expiresAt",
          cell: ({ row }) => {
            const val = row.original.expiresAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "expiresAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Expires At",
          },
        },
        {
          accessorKey: "ipAddress",
          cell: ({ row }) => {
            return row.original.ipAddress;
          },
          id: "ipAddress",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "IP Address",
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
          accessorKey: "token",
          cell: ({ row }) => {
            return row.original.token;
          },
          id: "token",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Token",
          },
        },
        {
          accessorKey: "userAgent",
          cell: ({ row }) => {
            return row.original.userAgent;
          },
          id: "userAgent",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "User Agent",
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
      deleteItem={async (id) => {
        await deleteSession(id);
      }}
      entityKey="sessions"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (session) => {
        await navigate({
          params: {
            sessionID: session.id,
          },
          to: `/sessions/$sessionID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={SessionForm}
    />
  );
}
