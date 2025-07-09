import type { MemberEntity } from '@archesai/schemas'

import {
  deleteMember,
  getFindManyMembersSuspenseQueryOptions
} from '@archesai/client'
import { MEMBER_ENTITY_KEY } from '@archesai/schemas'
import { CheckIcon, User, XIcon } from '@archesai/ui/components/custom/icons'
import { DataTable } from '@archesai/ui/components/datatable/data-table'
import { Badge } from '@archesai/ui/components/shadcn/badge'

import MemberForm from '#components/forms/member-form'

export default function MemberDataTable() {
  return (
    <DataTable<MemberEntity>
      columns={[
        {
          accessorKey: 'role',
          cell: ({ row }) => {
            return (
              <Badge
                className='capitalize'
                variant={'secondary'}
              >
                {row.original.role.toLowerCase()}
              </Badge>
            )
          }
        },
        {
          accessorKey: 'userId',
          cell: ({ row }) => {
            return <span className='font-medium'>{row.original.userId}</span>
          }
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return (
              <span className='max-w-[500px] truncate font-medium'>
                {row.original.createdAt ?
                  <CheckIcon className='text-primary' />
                : <XIcon className='text-red-950' />}
              </span>
            )
          }
        }
      ]}
      // content={() => (
      //   <div className='flex h-full w-full items-center justify-center'>
      //     <User
      //       className='opacity-30'
      //       size={100}
      //     />
      //   </div>
      // )}
      createForm={<MemberForm />}
      defaultView='table'
      deleteItem={async (id) => {
        await deleteMember(id)
      }}
      entityType={MEMBER_ENTITY_KEY}
      getEditFormFromItem={(member) => {
        return <MemberForm memberId={member.id} />
      }}
      handleSelect={() => {
        console.log('handleSelect')
      }}
      icon={<User />}
      useFindMany={getFindManyMembersSuspenseQueryOptions()}
    />
  )
}
