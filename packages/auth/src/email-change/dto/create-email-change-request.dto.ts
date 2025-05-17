import type { StaticDecode } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateEmailChangeRequestSchema = Type.Object({
  newEmail: Type.String({
    description: 'The e-mail to send the confirmation token to',
    format: 'email'
  }),
  userId: Type.String({
    description: 'The user ID of the user requesting the email change',
    format: 'uuid'
  })
})

export type CreateEmailChangeRequest = StaticDecode<
  typeof CreateEmailChangeRequestSchema
>
