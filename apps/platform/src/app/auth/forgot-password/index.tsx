import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'
import { createFileRoute, Link } from '@tanstack/react-router'

import { useRequestPasswordReset } from '@archesai/client'
import { AuthForm } from '@archesai/ui/components/custom/auth-form'

const ForgotPasswordSchema = Type.Object({
  email: Type.String({ format: 'email' })
})

export const Route = createFileRoute('/auth/forgot-password/')({
  component: ForgotPasswordPage
})

export default function ForgotPasswordPage() {
  const { mutate: requestPasswordReset } = useRequestPasswordReset()

  const onSubmit = (data: Static<typeof ForgotPasswordSchema>) => {
    requestPasswordReset(
      {
        data: {
          email: data.email
        }
      },
      {
        onError: (err) => {
          console.log(err)
        },
        onSuccess: (session) => {
          console.log('Request Password Reset successful:', session)
        }
      }
    )
  }

  return (
    <>
      <AuthForm
        description='Enter your email address to receive a password reset link'
        fields={[
          {
            defaultValue: '',
            label: 'Email',
            name: 'email',
            type: 'email',
            validationRule: ForgotPasswordSchema.properties.email
          }
        ]}
        onSubmit={(data) => {
          onSubmit(data as Static<typeof ForgotPasswordSchema>)
        }}
        title='Forgot Password'
      />
      <div className='text-center text-sm'>
        Remembered your password?{' '}
        <Link
          className='underline'
          to='/auth/login'
        >
          Login
        </Link>
      </div>
    </>
  )
}
