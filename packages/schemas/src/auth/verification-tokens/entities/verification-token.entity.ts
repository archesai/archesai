import type { Static, TObject, TOptional, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const VerificationTokenEntitySchema: TObject<{
  createdAt: TString
  expires: TString
  id: TString
  identifier: TString
  newEmail: TOptional<TString>
  token: TString
  updatedAt: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    expires: Type.String({
      description: 'The expiration date of the token'
    }),
    identifier: Type.String({
      description: 'The identifier associated with the token'
    }),
    newEmail: Type.Optional(
      Type.String({
        description: 'The new email address (if applicable)',
        format: 'email'
      })
    ),
    token: Type.String({
      description: 'The token string'
    })
  },
  {
    $id: 'VerificationTokenEntity',
    description: 'The verification token entity',
    title: 'Verification Token Entity'
  }
)

export type VerificationTokenEntity = Static<
  typeof VerificationTokenEntitySchema
>

export const VERIFICATION_TOKEN_ENTITY_KEY = 'verification-tokens'
