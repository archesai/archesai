import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'
import { RoleTypes } from '#enums/role'

export const MemberEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  id: z.ZodString
  organizationId: z.ZodString
  role: z.ZodEnum<{
    admin: 'admin'
    member: 'member'
    owner: 'owner'
  }>
  updatedAt: z.ZodString
  userId: z.ZodString
}> = BaseEntitySchema.extend({
  organizationId: z.string().describe('The organization name'),
  role: z.enum(RoleTypes).describe('The role of the member'),
  userId: z.string().describe('The user id')
}).meta({
  description: 'Schema for Member entity',
  id: 'MemberEntity'
})

export type MemberEntity = z.infer<typeof MemberEntitySchema>

export const MEMBER_ENTITY_KEY = 'members'
