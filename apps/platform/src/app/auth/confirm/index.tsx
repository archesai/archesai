import { Suspense } from 'react'
import { createFileRoute } from '@tanstack/react-router'

import { ConfirmationForm } from '@archesai/ui/components/custom/verification-token-confirmation-form'

type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export const Route = createFileRoute('/auth/confirm/')({
  component: ConfirmPage,
  validateSearch: (search) => {
    const { token, type } = search as {
      token: string
      type: '' | ActionType
    }
    return {
      token,
      type
    }
  }
})

export default function ConfirmPage() {
  const search = Route.useSearch()

  return (
    <>
      <Suspense fallback={<div>Loading...</div>}>
        <ConfirmationForm
          token={search.token}
          type={search.type || 'password-reset'}
        />
      </Suspense>
    </>
  )
}
