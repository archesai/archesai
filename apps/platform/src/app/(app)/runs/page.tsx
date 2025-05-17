import type { Metadata } from 'next'

import RunDataTable from '#components/datatables/run-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/runs')

export default function RunsPage() {
  return <RunDataTable />
}
