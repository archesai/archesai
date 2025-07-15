import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

export const CreateEmailVerificationDtoSchema: TObject<{
  email: TString
  userId: TString
}> = Type.Object({
  email: Type.String({
    description: 'The e-mail to send the confirmation token to',
    format: 'email'
  }),
  userId: Type.String({
    description: 'The user ID of the user requesting the email verification',
    format: 'uuid'
  })
})

export type CreateEmailVerificationDto = Static<
  typeof CreateEmailVerificationDtoSchema
>
