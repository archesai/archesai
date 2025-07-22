import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const ApiTokenEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  enabled: z.ZodBoolean
  expiresAt: z.ZodNullable<z.ZodString>
  id: z.ZodUUID
  key: z.ZodString
  lastRefill: z.ZodNullable<z.ZodString>
  lastRequest: z.ZodNullable<z.ZodString>
  metadata: z.ZodNullable<z.ZodString>
  name: z.ZodNullable<z.ZodString>
  permissions: z.ZodNullable<z.ZodString>
  prefix: z.ZodNullable<z.ZodString>
  rateLimitEnabled: z.ZodBoolean
  rateLimitMax: z.ZodNullable<z.ZodNumber>
  rateLimitTimeWindow: z.ZodNullable<z.ZodNumber>
  refillAmount: z.ZodNullable<z.ZodNumber>
  refillInterval: z.ZodNullable<z.ZodNumber>
  remaining: z.ZodNullable<z.ZodNumber>
  requestCount: z.ZodNumber
  start: z.ZodNullable<z.ZodString>
  updatedAt: z.ZodString
  userId: z.ZodString
}> = BaseEntitySchema.extend({
  enabled: z.boolean().describe('Whether the API token is enabled or not'),
  expiresAt: z
    .string()
    .nullable()
    .describe('The date and time when the API token expires'),
  key: z.string().describe('The API token key. This will only be shown once'),
  lastRefill: z
    .string()
    .nullable()
    .describe('The date and time when the API token was last refilled'),
  lastRequest: z
    .string()
    .nullable()
    .describe('The date and time when the API token was last used'),
  metadata: z
    .string()
    .nullable()
    .describe('The metadata for the API token, used for custom data'),
  name: z.string().nullable().describe('The name of the API token'),
  permissions: z.string().nullable().describe('The name of the API token'),
  prefix: z
    .string()
    .nullable()
    .describe('TThe prefix for the API token, used for routing requests'),
  rateLimitEnabled: z
    .boolean()
    .describe('Whether the API token has rate limiting enabled'),
  rateLimitMax: z
    .number()
    .nullable()
    .describe('The maximum number of requests allowed per time window'),
  rateLimitTimeWindow: z
    .number()
    .nullable()
    .describe('The time window in seconds for the rate limit'),
  refillAmount: z
    .number()
    .nullable()
    .describe('The amount of requests to refill the token with'),
  refillInterval: z
    .number()
    .nullable()
    .describe('The interval in seconds to refill the token'),
  remaining: z
    .number()
    .nullable()
    .describe('The number of requests remaining for the token'),
  requestCount: z
    .number()
    .describe('The number of requests made with the token'),
  start: z
    .string()
    .nullable()
    .describe('The number of requests remaining for the token'),
  userId: z.string().describe('The id of the user the token belongs to')
}).meta({
  description: 'Schema for API Token entity',
  id: 'ApiTokenEntity'
})

export type ApiTokenEntity = z.infer<typeof ApiTokenEntitySchema>

export const API_TOKEN_ENTITY_KEY = 'api-tokens'
