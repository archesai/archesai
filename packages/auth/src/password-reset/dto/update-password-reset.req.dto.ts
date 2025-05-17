import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const UpdatePasswordResetRequestSchema = Type.Object({
  newPassword: Type.String({
    description: 'The new password'
  }),
  token: Type.String({
    description: 'The password reset token'
  })
})

export type UpdatePasswordResetRequest = Static<
  typeof UpdatePasswordResetRequestSchema
>
