'use client'

import type { Static } from '@sinclair/typebox'

import { useEffect, useState } from 'react'
import { typeboxResolver } from '@hookform/resolvers/typebox'
import { Type } from '@sinclair/typebox'
import { useForm } from 'react-hook-form'

import {
  confirmEmailChange,
  confirmEmailVerification,
  confirmPasswordReset
} from '@archesai/client'
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
import { toast } from '@archesai/ui/components/shadcn/sonner'

// Define allowed action types
type ActionType = 'email-change' | 'email-verification' | 'password-reset'

const ConfirmationSchema = Type.Object({
  confirmPassword: Type.String({
    maxLength: 128,
    minLength: 8
  }),
  password: Type.String({
    maxLength: 128,
    minLength: 8
  })
})

export const ConfirmationForm = ({
  token,
  type
}: {
  token: string
  type: ActionType
}) => {
  const [operationSent, setOperationSent] = useState<boolean>(false)

  const form = useForm({
    defaultValues: {
      confirmPassword: '',
      password: ''
    },
    resolver: typeboxResolver(ConfirmationSchema)
  })

  const handleAction = async (token: string, type: ActionType) => {
    if (!token) {
      form.setError('root', {
        message: 'Invalid request. Missing parameters.'
      })
      return
    }
    if (operationSent) {
      return
    }
    setOperationSent(true)

    switch (type) {
      case 'email-change':
        try {
          await confirmEmailChange({
            newEmail: '',
            token,
            userId: ''
          })
          toast('Success', {
            description: 'Your email has been successfully updated!'
          })
        } catch (error) {
          console.error(error)
          form.setError('root', {
            message: 'Email change failed. Please try again.'
          })
        }
        break
      case 'email-verification':
        try {
          await confirmEmailVerification({
            token
          })
          toast('Your email has been successfully verified!')
        } catch (error) {
          console.error(error)
          form.setError('root', {
            message: 'Email verification failed. Please try again.'
          })
        }
        break
      case 'password-reset':
        break
      default:
        form.setError('root', { message: 'Unsupported action type.' })
        break
    }
  }

  useEffect(() => {
    if (type === 'email-change' || type === 'email-verification') {
      handleAction(token, type).catch((error: unknown) => {
        console.error(error)
        form.setError('root', {
          message: 'An unexpected error occurred. Please try again.'
        })
      })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const submitPasswordReset = async ({
    password
  }: Static<typeof ConfirmationSchema>) => {
    try {
      await confirmPasswordReset({
        newPassword: password,
        token
      })
      toast('Your password has been successfully reset!')
    } catch (error) {
      console.error(error)
      form.setError('root', {
        message: 'Password reset failed. Please try again.'
      })
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
      <p className='text-muted-foreground text-sm'>
        {form.formState.errors.root
          ? ''
          : type === 'password-reset'
            ? 'Please follow the instructions below.'
            : 'Verifying...'}
      </p>
    </Form>
  )
}
