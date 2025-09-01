import { z } from 'zod'

export const CreatePasswordResetDtoSchema: z.ZodObject<{
  email: z.ZodString
}> = z.object({
  email: z.string().describe('The e-mail to send the password reset token to')
})

export type CreatePasswordResetDto = z.infer<
  typeof CreatePasswordResetDtoSchema
>
