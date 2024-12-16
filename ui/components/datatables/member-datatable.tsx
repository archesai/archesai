'use client'
import { DataTable } from '@/components/datatables/datatable/data-table'
import MemberForm from '@/components/forms/member-form'
import { Badge } from '@/components/ui/badge'
import {
  MembersControllerFindAllPathParams,
  MembersControllerRemoveVariables,
  useMembersControllerFindAll,
  useMembersControllerRemove
} from '@/generated/archesApiComponents'
import { MemberEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import { CheckIcon, User, XIcon } from 'lucide-react'

export default function MemberDataTable() {
  const { defaultOrgname } = useAuth()

  return (
    <DataTable<
      MemberEntity,
      MembersControllerFindAllPathParams,
      MembersControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: 'role',
          cell: ({ row }) => {
            return (
              <Badge
                variant={'secondary'}
                className='capitalize'
              >
                {row.original.role.toLowerCase()}
              </Badge>
            )
          }
        },
        {
          accessorKey: 'inviteEmail',
          cell: ({ row }) => {
            return (
              <span className='font-medium'>{row.original.inviteEmail}</span>
            )
          }
        },
        {
          accessorKey: 'inviteAccepted',
          cell: ({ row }) => {
            return (
              <span className='max-w-[500px] truncate font-medium'>
                {row.original.inviteAccepted ? (
                  <CheckIcon className='text-primary' />
                ) : (
                  <XIcon className='text-red-950' />
                )}
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
      createForm={<MemberForm />}
      dataIcon={
        <User
          className='opacity-30'
          size={24}
        />
      }
      defaultView='table'
      filterField='inviteEmail'
      findAllPathParams={{
        orgname: defaultOrgname
      }}
      getDeleteVariablesFromItem={(member) => ({
        pathParams: {
          id: member.id,
          orgname: defaultOrgname
        }
      })}
      getEditFormFromItem={(member) => {
        return <MemberForm memberId={member.id} />
      }}
      handleSelect={() => {}}
      itemType='Member'
      useFindAll={useMembersControllerFindAll}
      useRemove={useMembersControllerRemove}
    />
  )
}
