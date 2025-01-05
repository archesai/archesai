import RunDataTable from '@/components/datatables/run-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/runs')

export default function RunsPage() {
  return <RunDataTable />
}
