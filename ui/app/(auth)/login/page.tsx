'use client'

import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/hooks/use-auth'
import { zodResolver } from '@hookform/resolvers/zod'
import Link from 'next/link'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import * as z from 'zod'

// Define schema using Zod for form validation
const schema = z.object({
  email: z.string().email({ message: 'Invalid email address' }),
  password: z.string().min(1, { message: 'Password is required' })
})

type LoginFormData = z.infer<typeof schema>

export default function LoginPage() {
  const { signInWithEmailAndPassword, signInWithGoogle } = useAuth()
  const [formError, setFormError] = useState<null | string>(null)

  const form = useForm<LoginFormData>({
    defaultValues: {
      email: '',
      password: ''
    },
    resolver: zodResolver(schema)
  })

  const onSubmit = async (data: LoginFormData) => {
    try {
      await signInWithEmailAndPassword(data.email, data.password)
      // Redirect handled by useEffect
    } catch (error: any) {
      console.error('Login error', error)
      // Enhanced error handling to capture specific error messages
      if (error?.message) {
        setFormError(error.message)
      } else {
        setFormError('An unexpected error occurred. Please try again.')
      }
    }
  }

  const handleGoogleSignIn = async () => {
    try {
      await signInWithGoogle()
      // Redirect handled by useEffect
    } catch (error: any) {
      console.error('Google Sign-In error', error)
      if (error?.message) {
        setFormError(error.message)
      } else {
        setFormError('An unexpected error occurred. Please try again.')
      }
    }
  }

  return (
    <div className='flex flex-col gap-2'>
      <div className='flex flex-col gap-2 text-center'>
        <h1 className='text-2xl font-semibold tracking-tight'>Login</h1>
        <p className='text-sm text-muted-foreground'>Enter your email and password to login to your account</p>
      </div>
      <div className='flex flex-col gap-2'>
        <Form {...form}>
          <form className='flex flex-col gap-2' noValidate onSubmit={form.handleSubmit(onSubmit)}>
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
                      aria-invalid={form.formState.errors.email ? 'true' : 'false'}
                    />
                  </FormControl>
                  <FormMessage>{form.formState.errors.email?.message}</FormMessage>
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
                      aria-invalid={form.formState.errors.password ? 'true' : 'false'}
                    />
                  </FormControl>
                  <FormMessage>{form.formState.errors.password?.message}</FormMessage>
                  {/* Forgot Password Link */}
                  <div className='text-right'>
                    <Link className='inline-block text-sm underline' href='/forgot-password'>
                      Forgot your password?
                    </Link>
                  </div>
                </FormItem>
              )}
            />

            {/* Display Form Error */}
            {formError && (
              <div className='text-center text-red-600' role='alert'>
                {formError}
              </div>
            )}

            {/* Submit Button */}
            <Button className='w-full' disabled={form.formState.isSubmitting} type='submit'>
              {form.formState.isSubmitting ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </Form>

        {/* Conditional Firebase Login Button */}
        {process.env.NEXT_PUBLIC_USE_FIREBASE === 'true' && (
          <Button
            className='w-full'
            disabled={form.formState.isSubmitting}
            onClick={handleGoogleSignIn}
            variant='outline'
          >
            {form.formState.isSubmitting ? 'Processing...' : 'Login with Google'}
          </Button>
        )}

        {/* Redirect to Register */}
        <div className='text-center text-sm'>
          Don&apos;t have an account?{' '}
          <Link className='underline' href='/register'>
            Sign up
          </Link>
        </div>
      </div>
    </div>
  )
}
