import { createFileRoute } from '@tanstack/react-router'

import ToolDataTable from '#components/datatables/tool-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/tools')

export const Route = createFileRoute('/(app)/tools/')({
  component: ToolsPage
})

export default function ToolsPage() {
  return <ToolDataTable />
}
