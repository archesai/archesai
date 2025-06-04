'use client'

import type { Static } from '@sinclair/typebox'

import Link from 'next/link'
// import { redirect } from 'next/navigation'
import { Type } from '@sinclair/typebox'

import { getSession, login } from '@archesai/client'

import { AuthForm } from '#components/auth-form'

const LoginSchema = Type.Object({
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export default function LoginPage() {
  const onSubmit = async (data: Static<typeof LoginSchema>) => {
    const response = await login(
      { email: data.email, password: data.password },
      {
        credentials: 'include'
      }
    )
    if (response.status === 401) {
      throw new Error(response.data.errors.map((e) => e.detail).join(', '))
    }
    const session = await getSession({
      credentials: 'include'
    })
    console.log('session', session)
    // redirect('/playground')
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
        onSubmit={(data) => onSubmit(data as Static<typeof LoginSchema>)}
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
