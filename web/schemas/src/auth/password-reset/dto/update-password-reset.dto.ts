import { z } from 'zod'

export const UpdatePasswordResetDtoSchema: z.ZodObject<{
  newPassword: z.ZodString
  token: z.ZodString
}> = z.object({
  newPassword: z.string().describe('The new password'),
  token: z.string().describe('The password reset token')
})

export type UpdatePasswordResetDto = z.infer<
  typeof UpdatePasswordResetDtoSchema
>
