'use client'

import type { Static } from '@sinclair/typebox'

import Link from 'next/link'
import { Type } from '@sinclair/typebox'

import { useLogin } from '@archesai/client'
import { AuthForm } from '@archesai/ui/components/custom/auth-form'

const LoginSchema = Type.Object({
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export default function LoginPage() {
  const { mutate: login } = useLogin()

  const onSubmit = (data: Static<typeof LoginSchema>) => {
    login(
      {
        data: {
          email: data.email,
          password: data.password
        }
      },
      {
        onError: (err) => {
          console.log(err)
        },
        onSuccess: (session) => {
          console.log('Login successful:', session)
        }
      }
    )
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
        onSubmit={(data) => {
          onSubmit(data as Static<typeof LoginSchema>)
        }}
        title='Login'
      />
      <div className='text-center text-sm'>
        Don&apos;t have an account?{' '}
        <Link
          className='underline'
          href='/auth/register'
        >
          Sign up
        </Link>
      </div>
    </>
  )
}
