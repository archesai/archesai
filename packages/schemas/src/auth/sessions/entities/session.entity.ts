import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const SessionEntitySchema: z.ZodObject<{
  activeOrganizationId: z.ZodString
  createdAt: z.ZodString
  expiresAt: z.ZodString
  id: z.ZodString
  ipAddress: z.ZodNullable<z.ZodString>
  token: z.ZodString
  updatedAt: z.ZodString
  userAgent: z.ZodNullable<z.ZodString>
  userId: z.ZodString
}> = BaseEntitySchema.extend({
  activeOrganizationId: z.string().describe('The active organization ID'),
  expiresAt: z.string().describe('The expiration date of the session'),
  ipAddress: z.string().nullable().describe('The IP address of the session'),
  token: z.string().describe('The session token'),
  userAgent: z.string().nullable().describe('The user agent of the session'),
  userId: z.string().describe('The ID of the user associated with the session')
}).meta({
  description: 'Schema for Session entity',
  id: 'SessionEntity'
})

export type SessionEntity = z.infer<typeof SessionEntitySchema>

export const SESSION_ENTITY_KEY = 'sessions'
