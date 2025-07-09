import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'
import { createFileRoute, Link, useRouter } from '@tanstack/react-router'

import { useLogin } from '@archesai/client'
import { AuthForm } from '@archesai/ui/components/custom/auth-form'

const LoginSchema = Type.Object({
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export const Route = createFileRoute('/auth/login/')({
  component: LoginPage
})

export default function LoginPage() {
  const router = useRouter()
  const { mutateAsync: login } = useLogin()

  const onSubmit = async (data: Static<typeof LoginSchema>) => {
    await login({
      data: {
        email: data.email,
        password: data.password
      }
    })
    await router.navigate({ to: '/chat' })
  }

  return (
    <>
      <AuthForm
        description='Enter your email and password to login'
        fields={[
          {
            defaultValue: '',
            label: 'Email',
            name: 'email',
            type: 'email',
            validationRule: LoginSchema.properties.email
          },
          {
            defaultValue: '',
            label: 'Password',
            name: 'password',
            type: 'password',
            validationRule: LoginSchema.properties.password
          }
        ]}
        onSubmit={async (data) => {
          await onSubmit(data as Static<typeof LoginSchema>)
        }}
        title='Login'
      />
      <div className='text-center text-sm'>
        Don&apos;t have an account?{' '}
        <Link
          className='underline'
          to='/auth/register'
        >
          Sign up
        </Link>
      </div>
    </>
  )
}
