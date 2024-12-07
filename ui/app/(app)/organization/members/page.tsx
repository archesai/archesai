import MemberDataTable from '@/components/datatables/member-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/organization/members')

export default function MembersPage() {
  return <MemberDataTable />
}
