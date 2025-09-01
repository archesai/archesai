import type { z } from 'zod'

import { CreateUserDtoSchema } from '#auth/users/dto/create-user.dto'

export const UpdateUserDtoSchema: z.ZodObject<{
  email: z.ZodOptional<z.ZodString>
  image: z.ZodOptional<z.ZodNullable<z.ZodString>>
}> = CreateUserDtoSchema.partial()

export type UpdateUserDto = z.infer<typeof UpdateUserDtoSchema>
