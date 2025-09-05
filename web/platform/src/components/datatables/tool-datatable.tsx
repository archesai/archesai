import type { JSX } from 'react'

import { Link, useNavigate } from '@tanstack/react-router'

import type { FindManyToolsParams, ToolEntity } from '@archesai/client'
import type { SearchQuery } from '@archesai/ui/types/entities'

import {
  deleteTool,
  getFindManyToolsSuspenseQueryOptions
} from '@archesai/client'
import {
  CalendarIcon,
  PackageCheckIcon,
  TextIcon
} from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'
import { TOOL_ENTITY_KEY } from '@archesai/ui/lib/constants'

export default function ToolDataTable(): JSX.Element {
  const navigate = useNavigate()

  const getQueryOptions = (query: SearchQuery) => {
    const params: any =
      query.filter || query.page || query.sort ?
        {
          ...(query.filter && {
            filter: query.filter as unknown as FindManyToolsParams['filter']
          }),
          ...(query.page && { page: query.page }),
          ...(query.sort && { sort: query.sort as FindManyToolsParams['sort'] })
        }
      : undefined
    return getFindManyToolsSuspenseQueryOptions(params) as any
  }

  return (
    <DataTable<ToolEntity>
      columns={[
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <Link
                className='shrink truncate text-wrap text-blue-500 hover:underline md:text-sm'
                params={{
                  artifactId: row.original.id
                }}
                to={`/artifacts/$artifactId`}
              >
                {row.original.name}
              </Link>
            )
          },
          enableColumnFilter: true,
          id: 'name',
          meta: {
            filterVariant: 'text',
            icon: TextIcon,
            label: 'Name'
          }
        },
        {
          accessorKey: 'description',
          cell: ({ row }) => {
            return row.original.description || 'No Description'
          },
          enableColumnFilter: true,
          id: 'description',
          meta: {
            filterVariant: 'text',
            icon: TextIcon,
            label: 'Description'
          }
        },
        {
          accessorKey: 'inputMimeType',
          cell: ({ row }) => {
            return (
              <Badge variant={'secondary'}>{row.original.inputMimeType}</Badge>
            )
          },
          enableColumnFilter: true,
          id: 'inputMimeType',
          meta: {
            filterVariant: 'text',
            icon: TextIcon,
            label: 'Input'
          }
        },
        {
          accessorKey: 'outputMimeType',
          cell: ({ row }) => {
            return (
              <Badge variant={'secondary'}>{row.original.outputMimeType}</Badge>
            )
          },
          enableColumnFilter: true,
          id: 'outputMimeType',
          meta: {
            filterVariant: 'multiSelect',
            icon: TextIcon,
            label: 'Output'
          }
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          enableColumnFilter: true,
          id: 'createdAt',
          meta: {
            filterVariant: 'dateRange',
            icon: CalendarIcon,
            label: 'Created at'
          }
        }
      ]}
      deleteItem={async (id) => {
        await deleteTool(id)
      }}
      entityKey={TOOL_ENTITY_KEY}
      getQueryOptions={getQueryOptions}
      handleSelect={async (tool) => {
        await navigate({ to: `/tool/single?toolId=${tool.id}` })
      }}
      icon={<PackageCheckIcon />}
    />
  )
}
