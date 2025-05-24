import { RocketIcon } from 'lucide-react'
import { toast } from 'sonner'

import type { UserEntity } from '@archesai/domain'

import { requestEmailVerification } from '@archesai/client'

import { Alert, AlertTitle } from '#components/shadcn/alert'

export interface VerifyEmailAlertProps {
  user?: UserEntity
}

export function VerifyEmailAlert({ user }: VerifyEmailAlertProps) {
  const handleRequestEmailVerification = async () => {
    try {
      await requestEmailVerification()
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
