import { Suspense } from 'react'
import Link from 'next/link'

import { Button } from '@archesai/ui/components/shadcn/button'

import { ConfirmationForm } from './password-reset-form'

// Define allowed action types
type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export default async function ConfirmPage({
  searchParams
}: {
  searchParams: Promise<Record<string, string | undefined>>
}) {
  const { token = '', type = '' } = (await searchParams) as {
    token: string
    type: '' | ActionType
  }

  return (
    <div className='flex flex-col gap-2'>
      <div className='flex flex-col gap-2 text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>
          {type.split('-').join(' ')}
        </h1>
      </div>
      <div className='flex flex-col gap-2'>
        {/* Handle Password Reset Form */}
        <Suspense fallback={<div>Loading...</div>}>
          <ConfirmationForm
            token={token}
            type={type || 'password-reset'}
          />
        </Suspense>
        <Suspense fallback={<div>Loading...</div>}>
          <div className='text-center'>
            <Button asChild>
              <Link href='/playground'>Go to Home</Link>
            </Button>
          </div>
        </Suspense>
      </div>
    </div>
  )
}
