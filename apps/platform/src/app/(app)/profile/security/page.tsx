'use client'

import {
  useDeleteUser,
  useGetOneUser,
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

export default function ProfileSecuritySettingsPage() {
  const { data: userResponse } = useGetOneUser('me')
  const { isPending: deactivatePending, mutateAsync: deactivateAccount } =
    useDeleteUser()
  const {
    isPending: requestPasswordResetPending,
    mutateAsync: requestPasswordReset
  } = useRequestPasswordReset()

  if (!userResponse || userResponse.status !== 200) return null
  const user = userResponse.data.data

  return (
    <div className='flex flex-col gap-3'>
      <Card>
        <CardHeader>
          <CardTitle>Reset Password</CardTitle>
          <CardDescription>
            If you would like to change your password, please click the button
            below. It will send you an email with instructions on how to reset.
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
                    email: user.attributes.email
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
                  id: user.id
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
  )
}
