import type { Metadata } from 'next'

import LabelDataTable from '#components/datatables/label-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/labels')

export default function LabelsPage() {
  return <LabelDataTable />
}
