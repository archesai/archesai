import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const InvitationEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  email: z.ZodString
  expiresAt: z.ZodString
  id: z.ZodString
  inviterId: z.ZodString
  organizationId: z.ZodString
  role: z.ZodEnum<{
    admin: 'admin'
    member: 'member'
    owner: 'owner'
  }>
  status: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  email: z.string().describe('The email of the invitated user'),
  expiresAt: z
    .string()
    .describe('The date and time when the invitation expires'),
  inviterId: z.string().describe('The user id of the inviter'),
  organizationId: z
    .string()
    .describe('The name of the organization the token belongs to'),
  role: z.enum(RoleTypes).describe('The role of the invitation'),
  status: z
    .string()
    .describe('The status of the invitation, e.g., pending, accepted, declined')
}).meta({
  description: 'Schema for Invitation entity',
  id: 'InvitationEntity'
})

export type InvitationEntity = z.infer<typeof InvitationEntitySchema>

export const INVITATION_ENTITY_KEY = 'invitations'
