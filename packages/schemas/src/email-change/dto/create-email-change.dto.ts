import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateEmailChangeDtoSchema = Type.Object({
  newEmail: Type.String({
    description: 'The e-mail to send the confirmation token to',
    format: 'email'
  }),
  userId: Type.String({
    description: 'The user ID of the user requesting the email change',
    format: 'uuid'
  })
})

export type CreateEmailChangeDto = Static<typeof CreateEmailChangeDtoSchema>
