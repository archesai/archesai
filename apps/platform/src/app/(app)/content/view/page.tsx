import type { Metadata } from 'next'

import ContentDataTable from '#components/datatables/content-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/content/view')

export default function ContentPage() {
  return <ContentDataTable />
}
