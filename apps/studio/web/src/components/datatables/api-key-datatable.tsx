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
import APIKeyForm from "#components/forms/api-key-form";
import type { APIKey, PageQueryParameter } from "#lib/index";
import { deleteAPIKey, getListAPIKeysSuspenseQueryOptions } from "#lib/index";

export default function APIKeyDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListAPIKeysSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<APIKey>
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
          accessorKey: "keyHash",
          cell: ({ row }) => {
            return row.original.keyHash;
          },
          id: "keyHash",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Key Hash",
          },
        },
        {
          accessorKey: "lastUsedAt",
          cell: ({ row }) => {
            const val = row.original.lastUsedAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "lastUsedAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Last Used At",
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
          accessorKey: "prefix",
          cell: ({ row }) => {
            return row.original.prefix;
          },
          id: "prefix",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Prefix",
          },
        },
        {
          accessorKey: "rateLimit",
          cell: ({ row }) => {
            return row.original.rateLimit;
          },
          id: "rateLimit",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Rate Limit",
          },
        },
        {
          accessorKey: "scopes",
          cell: ({ row }) => {
            return row.original.scopes;
          },
          id: "scopes",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Scopes",
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
      createForm={APIKeyForm}
      deleteItem={async (id) => {
        await deleteAPIKey(id);
      }}
      entityKey="api_keys"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (apiKey) => {
        await navigate({
          params: {
            apiKeyID: apiKey.id,
          },
          to: `/api-keys/$apiKeyID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={APIKeyForm}
    />
  );
}
