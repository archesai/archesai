'use client'

import type { Static } from '@sinclair/typebox'

import Link from 'next/link'
import { typeboxResolver } from '@hookform/resolvers/typebox'
import { Type } from '@sinclair/typebox'
import { useForm } from 'react-hook-form'

import { requestPasswordReset } from '@archesai/client'
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

const ForgotPasswordSchema = Type.Object({
  email: Type.String({
    format: 'email'
  })
})

export default function ForgotPasswordPage() {
  const form = useForm({
    defaultValues: {
      email: ''
    },
    resolver: typeboxResolver(ForgotPasswordSchema)
  })

  const onSubmit = async (data: Static<typeof ForgotPasswordSchema>) => {
    try {
      await requestPasswordReset({
        email: data.email
      })
      form.setError('root', {
        message:
          'If an account with that email exists, a password reset link has been sent.'
      })
    } catch (error) {
      console.error(error)
      form.setError('root', {
        message: 'An unexpected error occurred. Please try again.'
      })
    }
  }

  return (
    <div className='flex flex-col gap-2'>
      <div className='flex flex-col gap-2 text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>
          Forgot Password
        </h1>
        <p className='text-md text-muted-foreground'>
          Enter your email address to receive a password reset link.
        </p>
      </div>
      <div className='flex flex-col gap-2'>
        {/* Display Error Message */}
        {form.formState.errors.root && (
          <div
            className='text-red-600'
            role='alert'
          >
            {form.formState.errors.email?.message}
          </div>
        )}

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
                  <FormLabel htmlFor='email'>Email Address</FormLabel>
                  <FormControl>
                    <Input
                      autoComplete='email'
                      id='email'
                      placeholder='you@example.com'
                      type='email'
                      {...field}
                      aria-invalid={
                        form.formState.errors.email ? 'true' : 'false'
                      }
                    />
                  </FormControl>
                  <FormMessage>
                    {form.formState.errors.root?.message}
                  </FormMessage>
                </FormItem>
              )}
            />

            {/* Submit Button */}
            <Button
              className='mt-5 w-full'
              disabled={form.formState.isSubmitting}
              type='submit'
            >
              {form.formState.isSubmitting ? 'Sending...' : 'Send Reset Link'}
            </Button>
          </form>
        </Form>

        {/* Redirect to Login */}
        <div className='text-center text-sm'>
          Remembered your password?{' '}
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
