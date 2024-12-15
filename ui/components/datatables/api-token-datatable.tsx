'use client'
import { DataTable } from '@/components/datatables/datatable/data-table'
import APITokenForm from '@/components/forms/api-token-form'
import { Badge } from '@/components/ui/badge'
import {
  ApiTokensControllerFindAllPathParams,
  ApiTokensControllerRemoveVariables,
  useApiTokensControllerFindAll,
  useApiTokensControllerRemove
} from '@/generated/archesApiComponents'
import { ApiTokenEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { format } from 'date-fns'
import { User } from 'lucide-react'
import { useRouter } from 'next/navigation'

export default function ApiTokenDataTable() {
  const { defaultOrgname } = useAuth()
  const router = useRouter()

  return (
    <DataTable<
      ApiTokenEntity,
      ApiTokensControllerFindAllPathParams,
      ApiTokensControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: 'role',
          cell: ({ row }) => {
            return <Badge>{row.original.role}</Badge>
          }
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
            return (
              <span className='font-light'>
                {format(new Date(row.original.createdAt), 'M/d/yy h:mm a')}
              </span>
            )
          }
        }
      ]}
      content={() => (
        <div className='flex h-full w-full items-center justify-center'>
          <User
            className='opacity-30'
            size={100}
          />
        </div>
      )}
      createForm={<APITokenForm />}
      dataIcon={
        <User
          className='opacity-30'
          size={24}
        />
      }
      defaultView='table'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(apiToken) => ({
        pathParams: {
          id: apiToken.id,
          orgname: defaultOrgname
        }
      })}
      getEditFormFromItem={(apiToken) => (
        <APITokenForm apiTokenId={apiToken.id} />
      )}
      handleSelect={(apiToken) => router.push(`/apiTokens/${apiToken.id}/chat`)}
      itemType='API token'
      useFindAll={useApiTokensControllerFindAll}
      useRemove={useApiTokensControllerRemove}
    />
  )
}
