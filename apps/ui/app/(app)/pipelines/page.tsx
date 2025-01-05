import PipelineDataTable from '@/components/datatables/pipeline-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/pipelines')

export default function PipelinesPage() {
  return <PipelineDataTable />
}
