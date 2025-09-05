import { toast } from 'sonner'

import type { UserEntity } from '#types/entities'

import { RocketIcon } from '#components/custom/icons'
import { Alert, AlertTitle } from '#components/shadcn/alert'

export interface VerifyEmailAlertProps {
  onRequestEmailVerification: () => Promise<void>
  user?: UserEntity
}

export function VerifyEmailAlert({
  onRequestEmailVerification,
  user
}: VerifyEmailAlertProps) {
  const handleRequestEmailVerification = async () => {
    try {
      await onRequestEmailVerification()
      toast('Email verification sent', {
        description: 'Please check your inbox for the verification email'
      })
    } catch (error) {
      toast('Error sending verification email', {
        description:
          error instanceof Error ? error.message : 'An error occurred'
      })
    }
  }

  if (!user || user.emailVerified) {
    return null
  }

  return (
    <Alert className='flex items-center rounded-none border-none bg-amber-700'>
      <RocketIcon
        className='h-5 w-5'
        color='white'
      />
      <AlertTitle className='flex items-center font-normal text-primary-foreground'>
        <span className='flex gap-1'>
          Please
          <div
            className='cursor-pointer font-semibold underline'
            onClick={handleRequestEmailVerification}
          >
            {' '}
            verify your email address{' '}
          </div>{' '}
          to continue using the app.
        </span>
      </AlertTitle>
    </Alert>
  )
}
