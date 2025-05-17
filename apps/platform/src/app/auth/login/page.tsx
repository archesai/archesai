'use client'

import type { Static } from '@sinclair/typebox'

import { useState } from 'react'
import Link from 'next/link'
import { typeboxResolver } from '@hookform/resolvers/typebox'
import { Type } from '@sinclair/typebox'
import { useForm } from 'react-hook-form'

import { Button } from '@archesai/ui/components/shadcn/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@archesai/ui/components/shadcn/form'
import { Input } from '@archesai/ui/components/shadcn/input'
import { useAuth } from '@archesai/ui/hooks/use-auth'

const LoginSchema = Type.Object({
  email: Type.String({
    format: 'email'
  }),
  password: Type.String({
    maxLength: 128,
    minLength: 8
  })
})

export default function LoginPage() {
  const { signInWithEmailAndPassword } = useAuth()
  const [formError, setFormError] = useState<null | string>(null)

  const form = useForm({
    defaultValues: {
      email: '',
      password: ''
    },
    resolver: typeboxResolver(LoginSchema)
  })

  const onSubmit = async (data: Static<typeof LoginSchema>) => {
    try {
      await signInWithEmailAndPassword(data.email, data.password)
    } catch (error: unknown) {
      console.error(error)
      setFormError('An unexpected error occurred. Please try again.')
    }
  }

  return (
    <div className='flex flex-col gap-2'>
      <div className='flex flex-col gap-2 text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>Login</h1>
        <p className='text-muted-foreground text-sm'>
          Enter your emails and password to login to your account
        </p>
      </div>
      <div className='flex flex-col gap-2'>
        <Form {...form}>
          <form
            className='flex flex-col gap-2'
            noValidate
            onSubmit={form.handleSubmit(onSubmit)}
          >
            {/* Email Field */}
            <FormField
              control={form.control}
              name='email'
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor='email'>Email</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete='email'
                      id='email'
                      placeholder='m@example.com'
                      type='email'
                      {...field}
                      aria-invalid={
                        form.formState.errors.email ? 'true' : 'false'
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.email?.message}
                  </FormMessage>
                </FormItem>
              )}
            />

            {/* Password Field */}
            <FormField
              control={form.control}
              name='password'
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor='password'>Password</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete='current-password'
                      id='password'
                      placeholder='Enter your password'
                      type='password'
                      {...field}
                      aria-invalid={
                        form.formState.errors.password ? 'true' : 'false'
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.password?.message}
                  </FormMessage>
                  {/* Forgot Password Link */}
                  <div className='text-right'>
                    <Link
                      className='inline-block text-sm underline'
                      href='/auth/forgot-password'
                    >
                      Forgot your password?
                    </Link>
                  </div>
                </FormItem>
              )}
            />

            {/* Display Form Error */}
            {formError && (
              <div
                className='text-center text-red-600'
                role='alert'
              >
                {formError}
              </div>
            )}

            {/* Submit Button */}
            <Button
              className='w-full'
              disabled={form.formState.isSubmitting}
              type='submit'
            >
              {form.formState.isSubmitting ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </Form>

        {/* Redirect to Register */}
        <div className='text-center text-sm'>
          Don&apos;t have an account?{' '}
          <Link
            className='underline'
            href='/auth/register'
          >
            Sign up
          </Link>
        </div>
      </div>
    </div>
  )
}
