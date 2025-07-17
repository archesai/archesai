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
  CardFooter,
  CardHeader,
  CardTitle
} from '@archesai/ui/components/shadcn/card'
import { Separator } from '@archesai/ui/components/shadcn/separator'

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
          <CardFooter>
            <Button
              disabled={requestPasswordResetPending}
              onClick={async () => {
                await requestPasswordReset({
                  data: {
                    email: sessionData.user.email
                  }
                })
              }}
              size={'sm'}
              type='submit'
            >
              {requestPasswordResetPending && (
                <LoaderIcon className='animate-spin' />
              )}
              Reset Password
            </Button>
          </CardFooter>
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
          <CardFooter>
            <Button
              disabled={deactivatePending}
              onClick={async () => {
                await deactivateAccount({
                  id: sessionData.user.id
                })
              }}
              size='sm'
              variant={'destructive'}
            >
              {deactivatePending && <LoaderPinwheel className='animate-spin' />}
              Delete Account
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  )
}
