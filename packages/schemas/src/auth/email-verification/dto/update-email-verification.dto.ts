import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const UpdateEmailVerificationDtoSchema: TObject<{
  token: TString
}> = Type.Object({
  token: Type.String({
    description: 'The password reset token'
  })
})

export type UpdateEmailVerificationDto = Static<
  typeof UpdateEmailVerificationDtoSchema
>
