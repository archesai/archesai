import { Link, useNavigate } from '@tanstack/react-router'

import type { ToolEntity } from '@archesai/schemas'

import { deleteTool, getFindManyToolsQueryOptions } from '@archesai/client'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'
import { ContentTypeToIcon } from '@archesai/ui/components/custom/content-type-to-icon'
import { PackageCheck, Text } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'

export default function ToolDataTable() {
  const navigate = useNavigate()

  return (
    <DataTable<ToolEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <div className='flex gap-2'>
                <Link
                  className='shrink truncate text-wrap text-primary hover:underline md:text-sm'
                  params={{
                    artifactId: row.original.id
                  }}
                  to={`/artifacts/$artifactId`}
                >
                  {row.original.name}
                </Link>
              </div>
            )
          },
          enableColumnFilter: true,
          meta: {
            filterVariant: 'text',
            icon: Text,
            label: 'Name'
          }
        },
        {
          accessorKey: 'description',
          cell: ({ row }) => {
            return <span>{row.original.description || 'No Description'}</span>
          },
          enableColumnFilter: true,
          enableHiding: false,
          meta: {
            filterVariant: 'text',
            icon: Text,
            label: 'Description'
          }
        },
        {
          accessorKey: 'inputMimeType',
          cell: ({ row }) => {
            return (
              <div className='flex items-center gap-2'>
                <ContentTypeToIcon contentType={row.original.inputMimeType} />
                {row.original.inputMimeType}
              </div>
            )
          },
          enableColumnFilter: true,
          enableHiding: false,
          meta: {
            filterVariant: 'text',
            icon: Text,
            label: 'Input'
          }
        },
        {
          accessorKey: 'outputMimeType',
          cell: ({ row }) => {
            return (
              <div className='flex items-center gap-2'>
                <ContentTypeToIcon contentType={row.original.outputMimeType} />
                {row.original.outputMimeType}
              </div>
            )
          },
          enableColumnFilter: true,
          enableHiding: false,
          meta: {
            filterVariant: 'multiSelect',
            icon: Text,
            label: 'Output'
          }
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
      entityKey={TOOL_ENTITY_KEY}
      handleSelect={async (tool) => {
        await navigate({ to: `/tool/single?toolId=${tool.id}` })
      }}
      icon={<PackageCheck />}
      useFindMany={getFindManyToolsQueryOptions}
    />
  )
}
