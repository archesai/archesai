import type { z } from 'zod'

import { UserEntitySchema } from '#auth/users/entities/user.entity'

export const CreateUserDtoSchema: z.ZodObject<{
  email: z.ZodString
  image: z.ZodNullable<z.ZodString>
}> = UserEntitySchema.pick({
  email: true,
  image: true
})

export type CreateUserDto = z.infer<typeof CreateUserDtoSchema>
