import { Type } from '@sinclair/typebox'
import { createFileRoute, Link } from '@tanstack/react-router'

import { useRequestPasswordReset } from '@archesai/client'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

const ForgotPasswordSchema = Type.Object({
  email: Type.String({ format: 'email' })
})

export const Route = createFileRoute('/auth/forgot-password/')({
  component: ForgotPasswordPage
})

export default function ForgotPasswordPage() {
  const { mutateAsync: requestPasswordReset } = useRequestPasswordReset()

  return (
    <>
      <GenericForm<
        typeof ForgotPasswordSchema.static,
        typeof ForgotPasswordSchema.static
      >
        description='Enter your email address to receive a password reset link'
        entityKey='auth'
        fields={[
          {
            defaultValue: '',
            label: 'Email',
            name: 'email',
            renderControl: (field) => (
              <Input
                {...field}
                type='email'
              />
            ),
            validationRule: ForgotPasswordSchema.properties.email
          }
        ]}
        isUpdateForm={false}
        onSubmitCreate={async (data) => {
          await requestPasswordReset({
            data: {
              email: data.email
            }
          })
        }}
        postContent={
          <div className='text-center text-sm'>
            Remembered your password?{' '}
            <Link
              className='underline'
              to='/auth/login'
            >
              Login
            </Link>
          </div>
        }
        showCard={true}
        title='Forgot Password'
      />
    </>
  )
}
