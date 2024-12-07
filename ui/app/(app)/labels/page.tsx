import LabelDataTable from '@/components/datatables/label-datatable'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/labels')

export default function LabelsPage() {
  return <LabelDataTable />
}
