import type { Metadata } from 'next'

import ToolDataTable from '#components/datatables/tool-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/tools')

export default function ToolsPage() {
  return <ToolDataTable />
}
