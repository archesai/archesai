import type { Metadata } from 'next'

import PipelineDataTable from '#components/datatables/pipeline-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/pipelines')

export default function PipelinesPage() {
  return <PipelineDataTable />
}
