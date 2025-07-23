import type { MemberEntity } from '@archesai/schemas'

import {
  deleteMember,
  getFindManyMembersSuspenseQueryOptions
} from '@archesai/client'
import { MEMBER_ENTITY_KEY } from '@archesai/schemas'
import { UserIcon } from '@archesai/ui/components/custom/icons'
import { Timestamp } from '@archesai/ui/components/custom/timestamp'
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
              <Badge variant={'secondary'}>
                {row.original.role.toLowerCase()}
              </Badge>
            )
          },
          id: 'role'
        },
        {
          accessorKey: 'userId',
          cell: ({ row }) => {
            return row.original.userId
          },
          id: 'userId'
        },
        {
          accessorKey: 'createdAt',
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: 'createdAt'
        }
      ]}
      createForm={MemberForm}
      deleteItem={async (id) => {
        await deleteMember(id)
      }}
      entityKey={MEMBER_ENTITY_KEY}
      getQueryOptions={getFindManyMembersSuspenseQueryOptions}
      handleSelect={() => {
        // Handle member selection if needed
      }}
      icon={<UserIcon />}
      updateForm={MemberForm}
    />
  )
}
