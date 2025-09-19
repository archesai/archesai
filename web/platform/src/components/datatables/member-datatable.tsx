import type {
  Member,
  MembersFilterParameter,
  MembersSortParameter,
  PageQueryParameter,
} from "@archesai/client";
import {
  deleteMember,
  getListMembersSuspenseQueryOptions,
  useGetSessionSuspense,
} from "@archesai/client";
import { Badge, Timestamp, UserIcon } from "@archesai/ui";
import { MEMBER_ENTITY_KEY } from "@archesai/ui/lib/constants";
import type { SearchQuery } from "@archesai/ui/types/entities";
import type { JSX } from "react";
import { DataTableContainer } from "#components/datatables/data-table-container";

import MemberForm from "#components/forms/member-form";

export default function MemberDataTable(): JSX.Element {
  const { data: sessionData } = useGetSessionSuspense("current");
  const organizationID = sessionData.data.organizationID;

  const getQueryOptions = (query: SearchQuery) => {
    return getListMembersSuspenseQueryOptions(organizationID, {
      filter: query.filter as unknown as MembersFilterParameter,
      page: query.page as PageQueryParameter,
      sort: query.sort as MembersSortParameter,
    });
  };

  return (
    <DataTableContainer<Member>
      columns={[
        {
          accessorKey: "role",
          cell: ({ row }) => {
            return (
              <Badge variant={"secondary"}>
                {row.original.role.toLowerCase()}
              </Badge>
            );
          },
          id: "role",
        },
        {
          accessorKey: "userID",
          cell: ({ row }) => {
            return row.original.userID;
          },
          id: "userID",
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return <Timestamp date={row.original.createdAt} />;
          },
          id: "createdAt",
        },
      ]}
      createForm={MemberForm}
      deleteItem={async (id) => {
        await deleteMember(organizationID, id);
      }}
      entityKey={MEMBER_ENTITY_KEY}
      // biome-ignore lint/suspicious/noExplicitAny: FIXME
      getQueryOptions={getQueryOptions as any}
      handleSelect={() => {
        // Handle member selection if needed
      }}
      icon={<UserIcon />}
      updateForm={MemberForm}
    />
  );
}
