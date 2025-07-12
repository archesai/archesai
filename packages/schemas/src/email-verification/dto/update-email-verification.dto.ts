import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const UpdateEmailVerificationDtoSchema = Type.Object({
  token: Type.String({
    description: 'The password reset token'
  })
})

export type UpdateEmailVerificationDto = Static<
  typeof UpdateEmailVerificationDtoSchema
>
