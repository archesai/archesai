'use client'

import { useRouter } from 'next/navigation'

import type { ApiTokenEntity } from '@archesai/domain'

import { useFindManyApiTokens } from '@archesai/client'
import { API_TOKEN_ENTITY_KEY } from '@archesai/domain'
import { User } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

import APITokenForm from '#components/forms/api-token-form'

export default function ApiTokenDataTable() {
  const router = useRouter()

  return (
    <DataTable<ApiTokenEntity>
      columns={[
        {
          accessorKey: 'role',
          cell: ({ row }) => (
            <Badge
              className='capitalize'
              variant={'secondary'}
            >
              {row.original.role.toLowerCase()}
            </Badge>
          )
        },
        {
          accessorKey: 'name',
          cell: ({ row }) => {
            return (
              <span className='max-w-[500px] truncate font-medium'>
                {row.original.name}
              </span>
            )
          }
        },
        {
          accessorKey: 'key',
          cell: ({ row }) => {
            return <span className='font-medium'>{row.original.key}</span>
          }
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          }
        }
      ]}
      createForm={<APITokenForm />}
      defaultView='table'
      entityType={API_TOKEN_ENTITY_KEY}
      getEditFormFromItem={(apiToken) => (
        <APITokenForm apiTokenId={apiToken.id} />
      )}
      grid={() => (
        <div className='flex h-full w-full items-center justify-center'>
          <User
            className='opacity-30'
            size={100}
          />
        </div>
      )}
      handleSelect={(apiToken) => {
        console.error('handleSelect', apiToken)
        router.push(`/organization/api-tokens`)
      }}
      icon={
        <User
          className='opacity-30'
          size={24}
        />
      }
      useFindMany={useFindManyApiTokens}
    />
  )
}
