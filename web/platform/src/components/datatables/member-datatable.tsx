import type { JSX } from "react"

import type {
  MemberEntity,
  MembersFilterParameter,
  MembersSortParameter,
  PageQueryParameter
} from "@archesai/client"
import type { SearchQuery } from "@archesai/ui/types/entities"

import {
  deleteMember,
  getFindManyMembersSuspenseQueryOptions,
  useGetOneSessionSuspense
} from "@archesai/client"
import { UserIcon } from "@archesai/ui/components/custom/icons"
import { Timestamp } from "@archesai/ui/components/custom/timestamp"
import { DataTable } from "@archesai/ui/components/datatable/data-table"
import { Badge } from "@archesai/ui/components/shadcn/badge"
import { MEMBER_ENTITY_KEY } from "@archesai/ui/lib/constants"

import MemberForm from "#components/forms/member-form"

export default function MemberDataTable(): JSX.Element {
  const { data: sessionData } = useGetOneSessionSuspense("current")
  const organizationId = sessionData.data.activeOrganizationId

  const getQueryOptions = (query: SearchQuery) => {
    return getFindManyMembersSuspenseQueryOptions(organizationId, {
      filter: query.filter as unknown as MembersFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as MembersSortParameter
    })
  }

  return (
    <DataTable<MemberEntity>
      columns={[
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return (
              <Badge variant={"secondary"}>
                {row.original.role.toLowerCase()}
              </Badge>
            )
          },
          id: "role"
        },
        {
          accessorKey: "userId",
          cell: ({ row }) => {
            return row.original.userId
          },
          id: "userId"
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />
          },
          id: "createdAt"
        }
      ]}
      createForm={MemberForm}
      deleteItem={async (id) => {
        await deleteMember(organizationId, id)
      }}
      entityKey={MEMBER_ENTITY_KEY}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-unsafe-assignment
      getQueryOptions={getQueryOptions as any}
      handleSelect={() => {
        // Handle member selection if needed
      }}
      icon={<UserIcon />}
      updateForm={MemberForm}
    />
  )
}
