import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const UpdatePasswordResetDtoSchema: TObject<{
  newPassword: TString
  token: TString
}> = Type.Object({
  newPassword: Type.String({
    description: 'The new password'
  }),
  token: Type.String({
    description: 'The password reset token'
  })
})

export type UpdatePasswordResetDto = Static<typeof UpdatePasswordResetDtoSchema>
