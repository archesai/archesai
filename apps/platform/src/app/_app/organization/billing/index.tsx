import { createFileRoute } from '@tanstack/react-router'

import BillingPageContent from '#app/_app/organization/billing/content'

export const Route = createFileRoute('/_app/organization/billing/')({
  component: BillingPage
})

export default function BillingPage() {
  return <BillingPageContent />
}
