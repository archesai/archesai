import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const VerificationEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  expiresAt: z.ZodString
  id: z.ZodString
  identifier: z.ZodString
  updatedAt: z.ZodString
  value: z.ZodString
}> = BaseEntitySchema.extend({
  expiresAt: z.string().describe('The expiration date of the token'),
  identifier: z.string().describe('The identifier associated with the token'),
  value: z.string().describe('The token string')
}).meta({
  description: 'Schema for Verification Token entity',
  id: 'VerificationTokenEntity'
})

export type VerificationEntity = z.infer<typeof VerificationEntitySchema>

export const VERIFICATION_TOKEN_ENTITY_KEY = 'verification-tokens'
