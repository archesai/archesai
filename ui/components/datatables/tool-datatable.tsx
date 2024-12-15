'use client'
import { DataTable } from '@/components/datatables/datatable/data-table'
import { Badge } from '@/components/ui/badge'
import {
  ToolsControllerFindAllPathParams,
  ToolsControllerRemoveVariables,
  useToolsControllerFindAll,
  useToolsControllerRemove
} from '@/generated/archesApiComponents'
import { ToolEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { format } from 'date-fns'
import { PackageCheck } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

export default function ToolDataTable() {
  const router = useRouter()
  const { defaultOrgname } = useAuth()

  return (
    <DataTable<
      ToolEntity,
      ToolsControllerFindAllPathParams,
      ToolsControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <Link
                  className='shrink truncate text-wrap text-blue-600 underline md:text-sm'
                  href={`/playground?selectedTool=${JSON.stringify(row.original)}`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          }
        },
        {
          accessorKey: 'description',
          cell: ({ row }) => {
            return <span>{row.original.description || 'No Description'}</span>
          },
          enableHiding: false
        },
        {
          accessorKey: 'inputType',
          cell: ({ row }) => {
            return <Badge>{row.original.inputType}</Badge>
          },
          enableHiding: false
        },
        {
          accessorKey: 'outputType',
          cell: ({ row }) => {
            return <Badge>{row.original.outputType}</Badge>
          },
          enableHiding: false
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return (
              <span className='font-light'>
                {format(new Date(row.original.createdAt), 'M/d/yy h:mm a')}
              </span>
            )
          }
        }
      ]}
      dataIcon={<PackageCheck />}
      defaultView='table'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(tool) => ({
        pathParams: {
          id: tool.id,
          orgname: defaultOrgname
        }
      })}
      handleSelect={(tool) => router.push(`/tool/single?toolId=${tool.id}`)}
      itemType='tool'
      useFindAll={useToolsControllerFindAll}
      useRemove={useToolsControllerRemove}
    />
  )
}
