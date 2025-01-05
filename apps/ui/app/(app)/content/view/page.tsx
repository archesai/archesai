import ContentDataTable from '@/components/datatables/content-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/content/view')

export default function ContentPage() {
  return <ContentDataTable />
}
