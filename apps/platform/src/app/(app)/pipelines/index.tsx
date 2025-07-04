import { createFileRoute } from '@tanstack/react-router'

import PipelineDataTable from '#components/datatables/pipeline-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/pipelines')

export const Route = createFileRoute('/(app)/pipelines/')({
  component: PipelinesPage
})

export default function PipelinesPage() {
  return <PipelineDataTable />
}
