import type { StaticDecode } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const UpdateEmailVerificationRequestSchema = Type.Object({
  token: Type.String({
    description: 'The password reset token'
  })
})

export type UpdateEmailVerificationRequest = StaticDecode<
  typeof UpdateEmailVerificationRequestSchema
>
