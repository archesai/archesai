import type { Static, TObject, TString } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const VerificationEntitySchema: TObject<{
  createdAt: TString
  expiresAt: TString
  id: TString
  identifier: TString
  updatedAt: TString
  value: TString
}> = Type.Object(
  {
    ...BaseEntitySchema.properties,
    expiresAt: Type.String({
      description: 'The expiration date of the token'
    }),
    identifier: Type.String({
      description: 'The identifier associated with the token'
    }),
    value: Type.String({
      description: 'The token string'
    })
  },
  {
    $id: 'VerificationEntity',
    description: 'The verification token entity',
    title: 'Verification Token Entity'
  }
)

export type VerificationEntity = Static<typeof VerificationEntitySchema>

export const VERIFICATION_TOKEN_ENTITY_KEY = 'verification-tokens'
