import { createFileRoute } from '@tanstack/react-router'

import ArtifactDataTable from '#components/datatables/artifact-datatable'
import { getRouteMeta } from '#lib/site-utils'

export const metadata = getRouteMeta('/artifacts')

export const Route = createFileRoute('/_app/artifacts/')({
  component: ContentPage
})

export default function ContentPage() {
  return <ArtifactDataTable />
}
