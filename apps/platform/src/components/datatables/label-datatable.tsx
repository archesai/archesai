'use client'

import { useRouter } from 'next/navigation'

import type { LabelEntity } from '@archesai/domain'

import { deleteLabel, useFindManyLabels } from '@archesai/client'
import { LABEL_ENTITY_KEY } from '@archesai/domain'
import { ListMinus } from '@archesai/ui/components/custom/icons'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

import LabelForm from '#components/forms/label-form'

export default function LabelDataTable() {
  const router = useRouter()

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
                  onClick={() => {
                    router.push(`/chatbots/chat?labelId=${row.original.id}`)
                  }}
                >
                  <Badge variant={'secondary'}>{row.original.name}</Badge>
                </span>
              </div>
            )
          }
        }
      ]}
      createForm={<LabelForm />}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteLabel(id)
      }}
      entityType={LABEL_ENTITY_KEY}
      getEditFormFromItem={(label) => <LabelForm labelId={label.id} />}
      handleSelect={(chatbot) => {
        router.push(`/chatbots/chat?labelId=${chatbot.id}`)
      }}
      icon={<ListMinus />}
      useFindMany={useFindManyLabels}
    />
  )
}
