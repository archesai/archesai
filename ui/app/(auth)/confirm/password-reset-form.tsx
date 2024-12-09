'use client'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  fetchAuthControllerEmailChangeConfirm as changeEmail,
  fetchAuthControllerEmailVerificationConfirm as verifyEmail,
  fetchAuthControllerPasswordResetConfirm as resetPassword
} from '@/generated/archesApiComponents'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import * as z from 'zod'

// Define allowed action types
type ActionType = 'email-change' | 'email-verification' | 'password-reset'

// Define schemas for different actions
const passwordResetSchema = z
  .object({
    confirmPassword: z
      .string()
      .min(8, { message: 'Please confirm your password' }),
    password: z
      .string()
      .min(8, { message: 'Password must be at least 8 characters' })
      .regex(/^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/, {
        message:
          'Password must contain at least one letter, one number, and one special character'
      })
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword']
  })

type PasswordResetFormData = z.infer<typeof passwordResetSchema>

export const PasswordResetForm = ({
  token,
  type
}: {
  token: string
  type: ActionType
}) => {
  const [message, setMessage] = useState<string>('')
  const [error, setError] = useState<string>('')
  const [operationSent, setOperationSent] = useState<boolean>(false)

  const form = useForm<PasswordResetFormData>({
    defaultValues: {
      confirmPassword: '',
      password: ''
    },
    resolver: zodResolver(passwordResetSchema)
  })

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const handleAction = async (token: string, type: ActionType) => {
    if (!type || !token) {
      setError('Invalid request. Missing parameters.')
      return
    }
    if (operationSent) {
      return
    }
    setOperationSent(true)

    switch (type) {
      case 'email-change':
        try {
          await changeEmail({
            body: {
              token
            }
          })
          setMessage('Your email has been successfully updated!')
        } catch (err: any) {
          console.error(err)
          setError(
            err?.response?.data?.message ||
              'Email change failed. Please try again.'
          )
        }
        break
      case 'email-verification':
        try {
          await verifyEmail({
            body: {
              token
            }
          })
          setMessage('Your email has been successfully verified!')
        } catch (err: any) {
          console.error(err)
          // setError(
          //   err?.response?.data?.message || "Email verification failed."
          // );
        }
        break
      case 'password-reset':
        // Do nothing
        break

      default:
        setError('Unsupported action type.')
        break
    }
  }

  const submitPasswordReset = async ({ password }: PasswordResetFormData) => {
    try {
      await resetPassword({
        body: {
          newPassword: password,
          token
        }
      })
      setMessage('Your password has been successfully reset!')
    } catch (err: any) {
      console.error('Password reset error:', err)
      setError(
        err?.response?.data?.message ||
          'Password reset failed. Please try again.'
      )
    }
  }

  return (
    <Form {...form}>
      <form
        className='flex flex-col gap-2'
        noValidate
        onSubmit={form.handleSubmit(submitPasswordReset)}
      >
        {/* New Password Field */}
        <FormField
          control={form.control}
          name='password'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='password'>New Password</FormLabel>
              <FormControl>
                <Input
                  autoComplete='new-password'
                  id='password'
                  placeholder='Enter your new password'
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

        {/* Confirm New Password Field */}
        <FormField
          control={form.control}
          name='confirmPassword'
          render={({ field }) => (
            <FormItem>
              <FormLabel htmlFor='confirmPassword'>
                Confirm New Password
              </FormLabel>
              <FormControl>
                <Input
                  autoComplete='new-password'
                  id='confirmPassword'
                  placeholder='Confirm your new password'
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

        {/* Submit Button */}
        <Button
          className='w-full'
          disabled={form.formState.isSubmitting}
          type='submit'
        >
          {form.formState.isSubmitting
            ? 'Resetting Password...'
            : 'Reset Password'}
        </Button>
      </form>
      <p className='text-sm text-muted-foreground'>
        {message ||
          (error
            ? ''
            : type === 'password-reset'
              ? 'Please follow the instructions below.'
              : 'Verifying...')}
      </p>
    </Form>
  )
}
