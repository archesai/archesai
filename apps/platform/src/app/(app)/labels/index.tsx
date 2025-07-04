import { createFileRoute } from '@tanstack/react-router'

import LabelDataTable from '#components/datatables/label-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/labels')

export const Route = createFileRoute('/(app)/labels/')({
  component: LabelsPage
})

export default function LabelsPage() {
  return <LabelDataTable />
}
