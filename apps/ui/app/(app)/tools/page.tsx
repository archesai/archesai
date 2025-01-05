import ToolDataTable from '@/components/datatables/tool-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/tools')

export default function ToolsPage() {
  return <ToolDataTable />
}
