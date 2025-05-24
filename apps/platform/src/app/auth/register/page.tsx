'use client'

import type { Static } from '@sinclair/typebox'

import { useState } from 'react'
import Link from 'next/link'
import { typeboxResolver } from '@hookform/resolvers/typebox'
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

import { RegisterSchema } from '#components/forms/register-form'

export default function RegisterPage() {
  const { registerWithEmailAndPassword } = useAuth()
  const [error, setError] = useState<null | string>(null)

  const form = useForm({
    defaultValues: {
      confirmPassword: '',
      email: '',
      password: ''
    },
    resolver: typeboxResolver(RegisterSchema)
  })

  const onSubmit = async (data: Static<typeof RegisterSchema>) => {
    try {
      await registerWithEmailAndPassword(data.email, data.password)
    } catch (error) {
      console.error(error)
      if (error instanceof Error) {
        setError(error.message)
      } else {
        setError('An unexpected error occurred. Please try again.')
      }
    }
  }

  return (
    <div className='flex flex-col gap-2'>
      <div className='flex flex-col gap-2 text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>Register</h1>
        <p className='text-sm text-muted-foreground'>
          Create your account by entering your email and password
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
                      autoComplete='new-password'
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
                </FormItem>
              )}
            />

            {/* Confirm Password Field */}
            <FormField
              control={form.control}
              name='confirmPassword'
              render={({ field }) => (
                <FormItem>
                  <FormLabel htmlFor='confirmPassword'>
                    Confirm Password
                  </FormLabel>
                  <FormControl>
                    <Input
                      autoComplete='new-password'
                      id='confirmPassword'
                      placeholder='Confirm your password'
                      type='password'
                      {...field}
                      aria-invalid={
                        form.formState.errors.confirmPassword ? 'true' : 'false'
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.confirmPassword?.message}
                  </FormMessage>
                </FormItem>
              )}
            />

            {/* Display Error Message */}
            {error && (
              <div
                className='text-center text-red-600'
                role='alert'
              >
                {error}
              </div>
            )}

            {/* Submit Button */}
            <Button
              className='mt-5 w-full'
              disabled={form.formState.isSubmitting}
              type='submit'
            >
              {form.formState.isSubmitting ? 'Registering...' : 'Register'}
            </Button>
          </form>
        </Form>

        {/* Redirect to Login */}
        <div className='text-center text-sm'>
          Already have an account?{' '}
          <Link
            className='underline'
            href='/auth/login'
          >
            Login
          </Link>
        </div>
      </div>
    </div>
  )
}
