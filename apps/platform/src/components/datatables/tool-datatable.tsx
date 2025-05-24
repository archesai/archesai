'use client'

import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { PackageCheck } from 'lucide-react'

import type { ToolEntity } from '@archesai/domain'

import { deleteTool, useFindManyTools } from '@archesai/client'
import { TOOL_ENTITY_KEY } from '@archesai/domain'
import { ContentTypeToIcon } from '@archesai/ui/components/custom/content-type-to-icon'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'

export default function ToolDataTable() {
  const router = useRouter()

  return (
    <DataTable<ToolEntity>
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
            return (
              <div className='flex h-full items-center justify-center'>
                <ContentTypeToIcon
                  contentType={row.original.inputType.toLowerCase()}
                />
              </div>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'outputType',
          cell: ({ row }) => {
            return (
              <div className='flex h-full items-center justify-center'>
                <ContentTypeToIcon
                  contentType={row.original.outputType.toLowerCase()}
                />
              </div>
            )
          },
          enableHiding: false
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          }
        }
      ]}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteTool(id)
      }}
      entityType={TOOL_ENTITY_KEY}
      handleSelect={(tool) => {
        router.push(`/tool/single?toolId=${tool.id}`)
      }}
      icon={<PackageCheck />}
      useFindMany={useFindManyTools}
    />
  )
}
