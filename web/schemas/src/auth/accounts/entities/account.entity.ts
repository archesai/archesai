import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const AccountEntitySchema: z.ZodObject<{
  accessToken: z.ZodNullable<z.ZodString>
  accessTokenExpiresAt: z.ZodNullable<z.ZodString>
  accountId: z.ZodString
  createdAt: z.ZodString
  id: z.ZodUUID
  idToken: z.ZodNullable<z.ZodString>
  password: z.ZodNullable<z.ZodString>
  providerId: z.ZodString
  refreshToken: z.ZodNullable<z.ZodString>
  refreshTokenExpiresAt: z.ZodNullable<z.ZodString>
  scope: z.ZodNullable<z.ZodString>
  updatedAt: z.ZodString
  userId: z.ZodString
}> = BaseEntitySchema.extend({
  accessToken: z.string().nullable().describe('The access token'),
  accessTokenExpiresAt: z.string().nullable().describe('The expiration date'),
  accountId: z.string().describe('The unique identifier for the account'),
  idToken: z.string().nullable().describe('The ID token'),
  password: z
    .string()
    .nullable()
    .describe('The hashed password for local authentication'),
  providerId: z
    .string()
    .describe('The provider ID associated with the auth provider'),
  refreshToken: z.string().nullable().describe('The refresh token'),
  refreshTokenExpiresAt: z
    .string()
    .nullable()
    .describe('The refresh token expiration date'),
  scope: z.string().nullable().describe('The scope of the access token'),
  userId: z.string().describe('The user ID associated with the auth provider')
}).meta({
  description: 'Schema for Account entity',
  id: 'AccountEntity'
})

export type AccountEntity = z.infer<typeof AccountEntitySchema>

export const ACCOUNT_ENTITY_KEY = 'accounts'
