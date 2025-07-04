import { Suspense } from 'react'

import { ConfirmationForm } from '@archesai/ui/components/custom/verification-token-confirmation-form'

type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export default async function ConfirmPage({
  searchParams
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>
}) {
  const { token = '', type = '' } = (await searchParams) as {
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
