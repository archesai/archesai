import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const VerificationTokenEntitySchema = Type.Object(
  {
    ...BaseEntitySchema.properties,
    expires: Type.String({
      description: 'The expiration date of the token',
      format: 'date-time'
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
