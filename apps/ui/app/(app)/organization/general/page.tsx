import OrganizationForm from '@/components/forms/organization-form'
import { getMetadata } from '@/config/site'
import { Metadata } from 'next'

export const metadata: Metadata = getMetadata('/organization/general')

export default function OrganizationPage() {
  return <OrganizationForm />
}
