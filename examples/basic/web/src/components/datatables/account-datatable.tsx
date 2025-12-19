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
import AccountForm from "#components/forms/account-form";
import type { Account, PageQueryParameter } from "#lib/index";
import { deleteAccount, getListAccountsSuspenseQueryOptions } from "#lib/index";

export default function AccountDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListAccountsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Account>
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
          accessorKey: "accessToken",
          cell: ({ row }) => {
            return row.original.accessToken;
          },
          id: "accessToken",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Access Token",
          },
        },
        {
          accessorKey: "accessTokenExpiresAt",
          cell: ({ row }) => {
            const val = row.original.accessTokenExpiresAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "accessTokenExpiresAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Access Token Expires At",
          },
        },
        {
          accessorKey: "accountIdentifier",
          cell: ({ row }) => {
            return row.original.accountIdentifier;
          },
          id: "accountIdentifier",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Account Identifier",
          },
        },
        {
          accessorKey: "idToken",
          cell: ({ row }) => {
            return row.original.idToken;
          },
          id: "idToken",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "ID Token",
          },
        },
        {
          accessorKey: "provider",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.provider}</Badge>;
          },
          id: "provider",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Provider",
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
          accessorKey: "refreshToken",
          cell: ({ row }) => {
            return row.original.refreshToken;
          },
          id: "refreshToken",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Refresh Token",
          },
        },
        {
          accessorKey: "refreshTokenExpiresAt",
          cell: ({ row }) => {
            const val = row.original.refreshTokenExpiresAt;
            return val ? <Timestamp date={val as string} /> : "-";
          },
          id: "refreshTokenExpiresAt",
          meta: {
            filterVariant: "date",
            icon: CalendarIcon,
            label: "Refresh Token Expires At",
          },
        },
        {
          accessorKey: "scope",
          cell: ({ row }) => {
            return row.original.scope;
          },
          id: "scope",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Scope",
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
        await deleteAccount(id);
      }}
      entityKey="accounts"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (account) => {
        await navigate({
          params: {
            accountID: account.id,
          },
          to: `/accounts/$accountID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={AccountForm}
    />
  );
}
