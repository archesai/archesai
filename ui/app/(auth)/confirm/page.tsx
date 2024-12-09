import { Button } from '@/components/ui/button'
import Link from 'next/link'
import { PasswordResetForm } from './password-reset-form'
import { Suspense } from 'react'

// Define allowed action types
type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export default async function ConfirmPage({
  searchParams
}: {
  searchParams: Promise<{ [key: string]: string | undefined }>
}) {
  const { type = '', token = '' } = (await searchParams) as {
    type: ActionType | ''
    token: string
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
          <PasswordResetForm
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
