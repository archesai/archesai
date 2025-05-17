import type { Metadata } from 'next'

import ApiTokenDataTable from '#components/datatables/api-token-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/organization/api-tokens')

export default function ApiTokensPage() {
  return <ApiTokenDataTable />
}
