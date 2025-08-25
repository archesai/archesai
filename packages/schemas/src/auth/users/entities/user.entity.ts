import { z } from 'zod'

import { BaseEntitySchema } from '#base/entities/base.entity'

export const UserEntitySchema: z.ZodObject<{
  createdAt: z.ZodString
  email: z.ZodString
  emailVerified: z.ZodBoolean
  id: z.ZodUUID
  image: z.ZodNullable<z.ZodString>
  name: z.ZodString
  updatedAt: z.ZodString
}> = BaseEntitySchema.extend({
  email: z.string().describe("The user's e-mail"),
  emailVerified: z
    .boolean()
    .describe("Whether or not the user's e-mail has been verified"),
  image: z.string().nullable().describe("The user's avatar image URL"),
  name: z.string().min(1).describe("The user's name")
}).meta({
  description: 'Schema for User entity',
  id: 'UserEntity'
})

export type UserEntity = z.infer<typeof UserEntitySchema>

export const USER_ENTITY_KEY = 'users'
