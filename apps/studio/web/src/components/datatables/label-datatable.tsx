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
import LabelForm from "#components/forms/label-form";
import type { Label, PageQueryParameter } from "#lib/index";
import { deleteLabel, getListLabelsSuspenseQueryOptions } from "#lib/index";

export default function LabelDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListLabelsSuspenseQueryOptions({
      page: query.page as PageQueryParameter,
    });
  };

  return (
    <DataTableContainer<Label>
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
      ]}
      createForm={LabelForm}
      deleteItem={async (id) => {
        await deleteLabel(id);
      }}
      entityKey="labels"
      // biome-ignore lint/suspicious/noExplicitAny: Query options type is compatible at runtime
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (label) => {
        await navigate({
          params: {
            labelID: label.id,
          },
          to: `/labels/$labelID`,
        });
      }}
      icon={<ListIcon />}
      updateForm={LabelForm}
    />
  );
}
