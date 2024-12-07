'use client'
// FIXME - remove use client
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import {
  useAuthControllerPasswordResetRequest,
  useUsersControllerDeactivate,
  useUsersControllerFindOne
} from '@/generated/archesApiComponents'
import { useAuth } from '@/hooks/use-auth'
import { useToast } from '@/hooks/use-toast'
import { ReloadIcon } from '@radix-ui/react-icons'

export default function ProfileSecuritySettingsPage() {
  const { data: user } = useUsersControllerFindOne({})
  const { toast } = useToast()
  const { logout } = useAuth()

  const { isPending: deactivatePending, mutateAsync: deactivateAccount } = useUsersControllerDeactivate()
  const { isPending: requestPasswordResetPending, mutateAsync: requestPasswordReset } =
    useAuthControllerPasswordResetRequest()

  return (
    <div className='flex flex-col gap-3'>
      <Card>
        <CardHeader>
          <CardTitle className='text-lg'>Reset Password</CardTitle>
          <CardDescription>
            If you would like to change your password, please click the button below. It will send you an email with
            instructions on how to reset.
          </CardDescription>
        </CardHeader>
        <Separator />
        <div className='flex justify-end rounded-lg bg-gray-50 p-4 py-2 dark:bg-black'>
          <Button
            disabled={requestPasswordResetPending}
            onClick={async () =>
              await requestPasswordReset(
                {
                  body: {
                    email: user?.email as string
                  }
                },
                {
                  onError: (err) => {
                    toast({
                      description: err?.message || 'An error occurred while trying to reset your password.',
                      title: 'Error'
                    })
                  },
                  onSuccess: () => {
                    toast({
                      description: 'We have sent you an email with instructions on how to reset your password.',
                      title: 'Email Sent'
                    })
                  }
                }
              )
            }
            size={'sm'}
          >
            {requestPasswordResetPending && <ReloadIcon className='h-5 w-5 animate-spin' />}
            <span> Reset Password</span>
          </Button>
        </div>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle className='text-lg'>Deactivate Account</CardTitle>
          <CardDescription>
            If you would like to deactivate your account, please click the button below. This action is irreversible.
          </CardDescription>
        </CardHeader>
        <Separator />
        <div className='flex justify-end rounded-lg bg-gray-50 p-4 py-2 dark:bg-black'>
          <Button
            disabled={deactivatePending}
            onClick={async () =>
              await deactivateAccount(
                {},
                {
                  onError: (err) => {
                    toast({
                      description: err?.message || 'An error occurred while trying to deactivate your account.',
                      title: 'Error'
                    })
                  },
                  onSuccess: async () => {
                    await logout()
                  }
                }
              )
            }
            size='sm'
            variant={'destructive'}
          >
            {deactivatePending && <ReloadIcon className='h-5 w-5 animate-spin' />}
            <span>Delete Account</span>
          </Button>
        </div>
      </Card>
    </div>
  )
}
