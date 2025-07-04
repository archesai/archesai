import { createFileRoute } from '@tanstack/react-router'

import ApiTokenDataTable from '#components/datatables/api-token-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/organization/api-tokens')

export const Route = createFileRoute('/_app/organization/api-tokens/')({
  component: ApiTokensPage
})

export default function ApiTokensPage() {
  return <ApiTokenDataTable />
}
