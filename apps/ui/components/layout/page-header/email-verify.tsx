import { Alert, AlertTitle } from '@/components/ui/alert'
import {
  useAuthControllerEmailVerificationRequest,
  useUsersControllerFindOne
} from '@/generated/archesApiComponents'
import { useToast } from '@/hooks/use-toast'
import { RocketIcon } from '@radix-ui/react-icons'

export function VerifyEmailAlert() {
  const { mutateAsync: requestEmailVerification } =
    useAuthControllerEmailVerificationRequest()
  const { toast } = useToast()
  const { data: user } = useUsersControllerFindOne({})

  if (!user || user?.emailVerified) return null
  return (
    <Alert className='flex items-center rounded-none border-none bg-amber-700'>
      <RocketIcon
        className='h-5 w-5'
        color='white'
      />
      <AlertTitle className='text-primary-foreground flex items-center font-normal'>
        <span className='flex gap-1'>
          Please
          <div
            className='cursor-pointer font-semibold underline'
            onClick={async () => {
              try {
                await requestEmailVerification({})
                toast({
                  description:
                    'Please check your inbox for the verification email',
                  title: 'Email verification sent'
                })
              } catch (error) {
                toast({
                  description: error as any,
                  title: 'Error sending verification email'
                })
              }
            }}
          >
            {' '}
            verify your email address{' '}
          </div>{' '}
          to continue using the app.
        </span>
      </AlertTitle>
    </Alert>
  )
}
