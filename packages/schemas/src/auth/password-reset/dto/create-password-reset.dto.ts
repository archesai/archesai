import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreatePasswordResetDtoSchema: TObject<{
  email: TString
}> = Type.Object({
  email: Type.String({
    description: 'The e-mail to send the password reset token to',
    format: 'email'
  })
})

export type CreatePasswordResetDto = Static<typeof CreatePasswordResetDtoSchema>
