import { createFileRoute, Link } from '@tanstack/react-router'

import type { CreatePasswordResetDto } from '@archesai/schemas'

import { useRequestPasswordReset } from '@archesai/client'
import { CreatePasswordResetDtoSchema } from '@archesai/schemas'
import { GenericForm } from '@archesai/ui/components/custom/generic-form'
import { Input } from '@archesai/ui/components/shadcn/input'

export const Route = createFileRoute('/auth/forgot-password/')({
  component: ForgotPasswordPage
})

export default function ForgotPasswordPage() {
  const { mutateAsync: requestPasswordReset } = useRequestPasswordReset()

  return (
    <>
      <GenericForm<CreatePasswordResetDto, CreatePasswordResetDto>
        description='Enter your email address to receive a password reset link'
        entityKey='auth'
        fields={[
          {
            label: 'Email',
            name: 'email',
            renderControl: (field) => (
              <Input
                {...field}
                type='email'
              />
            ),
            validationRule: CreatePasswordResetDtoSchema.shape.email
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
