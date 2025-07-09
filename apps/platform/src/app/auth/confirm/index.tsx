import { Suspense } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { ConfirmationForm } from '@archesai/ui/components/custom/verification-token-confirmation-form'

type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export const Route = createFileRoute('/auth/confirm/')({
  component: ConfirmPage
})

export default function ConfirmPage() {
  const search = Route.useSearch()
  const { token = '', type = '' } = search as {
    token: string
    type: '' | ActionType
  }

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <ConfirmationForm
          token={token}
          type={type || 'password-reset'}
        />
      </Suspense>
    </>
  )
}
