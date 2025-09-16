import type {
  Label,
  LabelsFilterParameter,
  LabelsSortParameter,
  PageQueryParameter,
} from "@archesai/client";
import {
  deleteLabel,
  getListLabelsSuspenseQueryOptions,
} from "@archesai/client";
import { Badge, ListIcon, Timestamp } from "@archesai/ui";
import { LABEL_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import { useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import { DataTableContainer } from "#components/datatables/data-table-container";

import LabelForm from "#components/forms/label-form";

export default function LabelDataTable(): JSX.Element {
  const navigate = useNavigate();

  const getQueryOptions = (query: SearchQuery) => {
    return getListLabelsSuspenseQueryOptions({
      filter: query.filter as unknown as LabelsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as LabelsSortParameter,
    });
  };

  return (
    <DataTableContainer<Label>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return <Badge variant={"secondary"}>{row.original.name}</Badge>;
          },
          id: "name",
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          id: "createdAt",
        },
        {
          accessorKey: "updatedAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.updatedAt} />;
          },
          id: "updatedAt",
        },
      ]}
      createForm={LabelForm}
      deleteItem={async (id) => {
        await deleteLabel(id);
      }}
      entityKey={LABEL_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (chatbot) => {
        await navigate({ to: `/chatbots/chat?labelId=${chatbot.id}` });
      }}
      icon={<ListIcon />}
      updateForm={LabelForm}
    />
  );
}
