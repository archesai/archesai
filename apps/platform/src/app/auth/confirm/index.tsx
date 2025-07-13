import { createFileRoute } from '@tanstack/react-router'

import { Type, Value } from '@archesai/schemas'
import { ConfirmationForm } from '@archesai/ui/components/custom/verification-token-confirmation-form'

export const Route = createFileRoute('/auth/confirm/')({
  component: ConfirmPage,
  validateSearch: (search) => {
    return Value.Parse(
      Type.Object({
        token: Type.String(),
        type: Type.Union(
          [
            Type.Literal('email-change'),
            Type.Literal('email-verification'),
            Type.Literal('password-reset')
          ],
          { errorMessage: 'Invalid action type' }
        )
      }),
      search
    )
  }
})

export default function ConfirmPage() {
  const search = Route.useSearch()

  return (
    <ConfirmationForm
      token={search.token}
      type={search.type}
    />
  )
}
