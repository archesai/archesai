import { createFileRoute } from '@tanstack/react-router'

import {
  useDeleteUser,
  useGetSessionSuspense,
  useRequestPasswordReset
} from '@archesai/client'
import {
  LoaderIcon,
  LoaderPinwheel
} from '@archesai/ui/components/custom/icons'
import { Button } from '@archesai/ui/components/shadcn/button'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'
import { Separator } from '@archesai/ui/components/shadcn/separator'
import { toast } from '@archesai/ui/components/shadcn/sonner'

import UserForm from '#components/forms/user-form'

export const Route = createFileRoute('/_app/profile/')({
  component: ProfileSecuritySettingsPage
})

export default function ProfileSecuritySettingsPage() {
  const { data: sessionData } = useGetSessionSuspense()
  const { isPending: deactivatePending, mutateAsync: deactivateAccount } =
    useDeleteUser()
  const {
    isPending: requestPasswordResetPending,
    mutateAsync: requestPasswordReset
  } = useRequestPasswordReset()

  return (
    <div className='flex flex-col gap-4'>
      <UserForm />
      <div className='grid grid-cols-1 gap-4 md:grid-cols-2'>
        <Card>
          <CardHeader>
            <CardTitle>Reset Password</CardTitle>
            <CardDescription>
              If you would like to change your password, please click the button
              below. It will send you an email with instructions on how to
              reset.
            </CardDescription>
          </CardHeader>
          <Separator />
          <div className='flex justify-end rounded-xl p-4 py-2'>
            <Button
              disabled={requestPasswordResetPending}
              onClick={async () => {
                await requestPasswordReset(
                  {
                    data: {
                      email: sessionData.user.email
                    }
                  },

                  {
                    onError: () => {
                      toast('Error', {
                        description:
                          'An error occurred while trying to reset your password.'
                      })
                    },
                    onSuccess: () => {
                      toast('Email Sent', {
                        description:
                          'We have sent you an email with instructions on how to reset your password.'
                      })
                    }
                  }
                )
              }}
              size={'sm'}
            >
              {requestPasswordResetPending && (
                <LoaderIcon className='h-5 w-5 animate-spin' />
              )}
              <span> Reset Password</span>
            </Button>
          </div>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Deactivate Account</CardTitle>
            <CardDescription>
              If you would like to deactivate your account, please click the
              button below. This action is irreversible.
            </CardDescription>
          </CardHeader>
          <Separator />
          <div className='flex justify-end rounded-xl p-4 py-2'>
            <Button
              disabled={deactivatePending}
              onClick={async () => {
                await deactivateAccount(
                  {
                    id: sessionData.user.id
                  },
                  {
                    onError: () => {
                      toast('Error', {
                        description:
                          'An error occurred while trying to deactivate your account.'
                      })
                    },
                    onSuccess: () => {
                      console.log('FIXME')
                    }
                  }
                )
              }}
              size='sm'
              variant={'destructive'}
            >
              {deactivatePending && (
                <LoaderPinwheel className='h-5 w-5 animate-spin' />
              )}
              <span>Delete Account</span>
            </Button>
          </div>
        </Card>
      </div>
    </div>
  )
}
