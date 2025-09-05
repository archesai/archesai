import type { JSX } from "react"

import { useNavigate } from "@tanstack/react-router"

import type {
  LabelEntity,
  LabelsFilterParameter,
  LabelsSortParameter,
  PageQueryParameter
} from "@archesai/client"
import type { SearchQuery } from "@archesai/ui/types/entities"

import {
  deleteLabel,
  getFindManyLabelsSuspenseQueryOptions
} from "@archesai/client"
import { ListIcon } from "@archesai/ui/components/custom/icons"
import { Timestamp } from "@archesai/ui/components/custom/timestamp"
import { DataTable } from "@archesai/ui/components/datatable/data-table"
import { Badge } from "@archesai/ui/components/shadcn/badge"
import { LABEL_ENTITY_KEY } from "@archesai/ui/lib/constants"

import LabelForm from "#components/forms/label-form"

export default function LabelDataTable(): JSX.Element {
  const navigate = useNavigate()

  const getQueryOptions = (query: SearchQuery) => {
    return getFindManyLabelsSuspenseQueryOptions({
      filter: query.filter as unknown as LabelsFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as LabelsSortParameter
    })
  }

  return (
    <DataTable<LabelEntity>
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return <Badge variant={"secondary"}>{row.original.name}</Badge>
          },
          id: "name"
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: "createdAt"
        },
        {
          accessorKey: "updatedAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.updatedAt} />
          },
          id: "updatedAt"
        }
      ]}
      createForm={LabelForm}
      deleteItem={async (id) => {
        await deleteLabel(id)
      }}
      entityKey={LABEL_ENTITY_KEY}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
      getQueryOptions={getQueryOptions as any}
      handleSelect={async (chatbot) => {
        await navigate({ to: `/chatbots/chat?labelId=${chatbot.id}` })
      }}
      icon={<ListIcon />}
      updateForm={LabelForm}
    />
  )
}
