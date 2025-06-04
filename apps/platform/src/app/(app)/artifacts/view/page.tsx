import type { Metadata } from 'next'

import ArtifactDataTable from '#components/datatables/artifact-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/artifacts/view')

export default function ContentPage() {
  return <ArtifactDataTable />
}
