'use client'
import { DataTable } from '@/components/datatables/datatable/data-table'
import { DataTableColumnHeader } from '@/components/datatables/datatable/data-table-column-header'
import LabelForm from '@/components/forms/label-form'
import {
  LabelsControllerFindAllPathParams,
  LabelsControllerRemoveVariables,
  useLabelsControllerFindAll,
  useLabelsControllerRemove
} from '@/generated/archesApiComponents'
import { LabelEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { ListMinus } from 'lucide-react'
import { useRouter } from 'next/navigation'

export default function LabelDataTable() {
  const { defaultOrgname } = useAuth()
  const router = useRouter()

  return (
    <DataTable<LabelEntity, LabelsControllerFindAllPathParams, LabelsControllerRemoveVariables>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <span
                  className='max-w-[500px] truncate font-medium'
                  onClick={() => router.push(`/chatbots/chat?labelId=${row.original.id}`)}
                >
                  {row.original.name}
                </span>
              </div>
            )
          },
          header: ({ column }) => <DataTableColumnHeader column={column} title='Name' />
        }
      ]}
      createForm={<LabelForm />}
      dataIcon={<ListMinus />}
      defaultView='table'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(label) => ({
        pathParams: {
          id: label.id,
          orgname: defaultOrgname
        }
      })}
      getEditFormFromItem={(label) => <LabelForm labelId={label.id} />}
      handleSelect={(chatbot) => router.push(`/chatbots/chat?labelId=${chatbot.id}`)}
      itemType='label'
      useFindAll={useLabelsControllerFindAll}
      useRemove={useLabelsControllerRemove}
    />
  )
}
