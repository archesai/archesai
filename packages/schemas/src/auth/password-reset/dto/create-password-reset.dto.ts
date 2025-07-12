import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreatePasswordResetDtoSchema = Type.Object({
  email: Type.String({
    description: 'The e-mail to send the password reset token to',
    format: 'email'
  })
})

export type CreatePasswordResetDto = Static<typeof CreatePasswordResetDtoSchema>
