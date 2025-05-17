import type { Metadata } from 'next'

import OrganizationForm from '#components/forms/organization-form'
import { getRouteMeta } from '#lib/site-utils'

export const metadata: Metadata = getRouteMeta('/organization/general')

export default function OrganizationPage() {
  return <OrganizationForm />
}
