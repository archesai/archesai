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
import OrganizationForm from "#components/forms/organization-form";
import type { Organization, PageQueryParameter } from "#lib/index";
import {
  deleteOrganization,
  getListOrganizationsSuspenseQueryOptions,
} from "#lib/index";

export default function OrganizationDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListOrganizationsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Organization>
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
          accessorKey: "billingEmail",
          cell: ({ row }) => {
            return row.original.billingEmail;
          },
          id: "billingEmail",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Billing Email",
          },
        },
        {
          accessorKey: "credits",
          cell: ({ row }) => {
            return row.original.credits;
          },
          id: "credits",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Credits",
          },
        },
        {
          accessorKey: "logo",
          cell: ({ row }) => {
            return row.original.logo;
          },
          id: "logo",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Logo",
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
          accessorKey: "plan",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.plan}</Badge>;
          },
          id: "plan",
          meta: {
            filterVariant: "multiSelect",
            icon: TextIcon,
            label: "Plan",
            options: [
              { label: "FREE", value: "FREE" },
              { label: "BASIC", value: "BASIC" },
              { label: "STANDARD", value: "STANDARD" },
              { label: "PREMIUM", value: "PREMIUM" },
              { label: "UNLIMITED", value: "UNLIMITED" },
            ],
          },
        },
        {
          accessorKey: "slug",
          cell: ({ row }) => {
            return row.original.slug;
          },
          id: "slug",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Slug",
          },
        },
        {
          accessorKey: "stripeCustomerIdentifier",
          cell: ({ row }) => {
            return row.original.stripeCustomerIdentifier;
          },
          id: "stripeCustomerIdentifier",
          meta: {
            filterVariant: "text",
            icon: TextIcon,
            label: "Stripe Customer Identifier",
          },
        },
      ]}
      createForm={OrganizationForm}
      deleteItem={async (id) => {
        await deleteOrganization(id);
      }}
      entityKey="organizations"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (organization) => {
        await navigate({
          params: {
            organizationID: organization.id,
          },
          to: `/organizations/$organizationID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={OrganizationForm}
    />
  );
}
