'use client'

import type { Static } from '@sinclair/typebox'

import Link from 'next/link'
import { redirect } from 'next/navigation'
import { Type } from '@sinclair/typebox'

import { register } from '@archesai/client'

import { AuthForm } from '#components/auth-form'

const RegisterSchema = Type.Object({
  confirmPassword: Type.String({ minLength: 8 }),
  email: Type.String({ format: 'email' }),
  password: Type.String({ minLength: 8 })
})

export default function RegisterPage() {
  const onSubmit = async (data: Static<typeof RegisterSchema>) => {
    const response = await register(
      {
        email: data.email,
        password: data.password
      },
      {
        credentials: 'include'
      }
    )
    if (response.status === 401) {
      throw new Error(response.data.errors.map((e) => e.detail).join(', '))
    }
    redirect('/playground')
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
        onSubmit={(data) => onSubmit(data as Static<typeof RegisterSchema>)}
        title='Register'
      />
      <div className='text-center text-sm'>
        Already have an account?{' '}
        <Link
          className='underline'
          href='/auth/login'
        >
          Login
        </Link>
      </div>
    </>
  )
}
