import { z } from 'zod'

import type { SessionEntity } from '#auth/sessions/entities/session.entity'
import type { UserEntity } from '#auth/users/entities/user.entity'

import { SessionEntitySchema } from '#auth/sessions/entities/session.entity'
import { UserEntitySchema } from '#auth/users/entities/user.entity'

export const BetterAuthSessionSchema: z.ZodObject<{
  session: z.ZodType<SessionEntity>
  user: z.ZodType<UserEntity>
}> = z.object({
  session: SessionEntitySchema,
  user: UserEntitySchema
})
