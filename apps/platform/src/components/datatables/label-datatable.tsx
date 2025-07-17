import { useNavigate } from '@tanstack/react-router'

import type { LabelEntity } from '@archesai/schemas'

import { deleteLabel, getFindManyLabelsQueryOptions } from '@archesai/client'
import { LABEL_ENTITY_KEY } from '@archesai/schemas'
import { ListMinus } from '@archesai/ui/components/custom/icons'
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
            return (
              <div className='flex gap-2'>
                <span
                  className='max-w-[500px] truncate font-medium'
                  onClick={async () => {
                    await navigate({
                      to: `/chatbots/chat?labelId=${row.original.id}`
                    })
                  }}
                >
                  <Badge variant={'secondary'}>{row.original.name}</Badge>
                </span>
              </div>
            )
          }
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
      icon={<ListMinus />}
      updateForm={LabelForm}
      useFindMany={getFindManyLabelsQueryOptions}
    />
  )
}
