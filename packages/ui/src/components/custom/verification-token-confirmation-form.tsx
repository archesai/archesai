import { useState } from 'react'

// import type { UpdatePasswordResetDto } from '@archesai/schemas'

// Define allowed action types
type ActionType = 'email-change' | 'email-verification' | 'password-reset'

export const ConfirmationForm = ({
  token,
  type
}: {
  token: string
  type: ActionType
}) => {
  const [operationSent, setOperationSent] = useState<boolean>(false)

  console.log('ConfirmationForm', { token, type })
  const _handleAction = (token: string, type: ActionType) => {
    if (operationSent) {
      return
    }
    setOperationSent(true)
    console.log('Handling action', { token, type })
    // switch (type) {
    //   case 'email-change':
    //     try {
    //       await confirmEmailChange({
    //         newEmail: '',
    //         token,
    //         userId: ''
    //       })
    //       toast('Success', {
    //         description: 'Your email has been successfully updated!'
    //       })
    //     } catch (error) {
    //       console.error(error)
    //       form.setError('root', {
    //         message: 'Email change failed. Please try again.'
    //       })
    //     }
    //     break
    //   case 'email-verification':
    //     try {
    //       await confirmEmailVerification({
    //         token
    //       })
    //       toast('Your email has been successfully verified!')
    //     } catch (error) {
    //       console.error(error)
    //       form.setError('root', {
    //         message: 'Email verification failed. Please try again.'
    //       })
    //     }
    //     break
    //   case 'password-reset':
    //     break
    //   default:
    //     form.setError('root', { message: 'Unsupported action type.' })
    //     break
    // }
    console.log(_handleAction)
  }

  return <></>
}
