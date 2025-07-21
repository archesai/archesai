import { useNavigate } from '@tanstack/react-router'

import type { LabelEntity } from '@archesai/schemas'

import { deleteLabel, getFindManyLabelsQueryOptions } from '@archesai/client'
import { LABEL_ENTITY_KEY } from '@archesai/schemas'
import { ListIcon } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

import LabelForm from '#components/forms/label-form'

export default function LabelDataTable() {
  const navigate = useNavigate()

  return (
    <DataTable<LabelEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return <Badge variant={'secondary'}>{row.original.name}</Badge>
          },
          id: 'name'
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: 'createdAt'
        },
        {
          accessorKey: 'updatedAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.updatedAt} />
          },
          id: 'updatedAt'
        }
      ]}
      createForm={LabelForm}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteLabel(id)
      }}
      entityKey={LABEL_ENTITY_KEY}
      handleSelect={async (chatbot) => {
        await navigate({ to: `/chatbots/chat?labelId=${chatbot.id}` })
      }}
      icon={<ListIcon />}
      updateForm={LabelForm}
      useFindMany={getFindManyLabelsQueryOptions}
    />
  )
}
