import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'
import { createFileRoute, Link, useRouter } from '@tanstack/react-router'

import { useRegister } from '@archesai/client'
import { AuthForm } from '@archesai/ui/components/custom/auth-form'

const RegisterSchema = Type.Object({
  confirmPassword: Type.String({ minLength: 8 }),
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export const Route = createFileRoute('/auth/register/')({
  component: RegisterPage
})

export default function RegisterPage() {
  const router = useRouter()
  const { mutateAsync: register } = useRegister()

  const onSubmit = async (data: Static<typeof RegisterSchema>) => {
    await register({
      data: {
        email: data.email,
        password: data.password
      }
    })
    await router.navigate({
      to: '/chat'
    })
  }

  return (
    <>
      <AuthForm
        description='Create your account by entering your email and password'
        fields={[
          {
            defaultValue: '',
            label: 'Email',
            name: 'email',
            type: 'email',
            validationRule: RegisterSchema.properties.email
          },
          {
            defaultValue: '',
            label: 'Password',
            name: 'password',
            type: 'password',
            validationRule: RegisterSchema.properties.password
          },
          {
            defaultValue: '',
            label: 'Confirm Password',
            name: 'confirmPassword',
            type: 'password',
            validationRule: Type.String({
              errorMessage: 'Passwords must match',
              minLength: 8
            })
          }
        ]}
        onSubmit={async (data) => {
          await onSubmit(data as Static<typeof RegisterSchema>)
        }}
        title='Register'
      />
      <div className='text-center text-sm'>
        Already have an account?{' '}
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
