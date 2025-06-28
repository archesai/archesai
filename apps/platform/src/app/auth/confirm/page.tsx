'use client'

import { Suspense } from 'react'
import Link from 'next/link'

import { ConfirmationForm } from '@archesai/ui/components/custom/verification-token-confirmation-form'
import { Button } from '@archesai/ui/components/shadcn/button'

type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export default function ConfirmPage({
  searchParams
}: {
  searchParams: Record<string, string | undefined>
}) {
  const { token = '', type = '' } = searchParams as {
    token: string
    type: '' | ActionType
  }

  const formatTitle = (actionType: string) => {
    return actionType
      .split('-')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ')
  }

  return (
    <>
      <div className='flex w-xs flex-col gap-2'>
        <div className='text-center'>
          <h1 className='text-2xl font-semibold tracking-tight'>
            {formatTitle(type || 'password-reset')}
          </h1>
        </div>
        <div className='flex flex-col gap-2'>
          <Suspense fallback={<div>Loading...</div>}>
            <ConfirmationForm
              token={token}
              type={type || 'password-reset'}
            />
          </Suspense>
        </div>
      </div>
      <div className='text-center text-sm'>
        <Button asChild>
          <Link href='/playground'>Go to Home</Link>
        </Button>
      </div>
    </>
  )
}
