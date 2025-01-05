import ApiTokenDataTable from '@/components/datatables/api-token-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/organization/api-tokens')

export default function ApiTokensPage() {
  return <ApiTokenDataTable />
}
