import type { Metadata } from 'next'

import MemberDataTable from '#components/datatables/member-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/organization/members')

export default function MembersPage() {
  return <MemberDataTable />
}
