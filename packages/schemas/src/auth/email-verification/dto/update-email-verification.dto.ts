import { z } from 'zod'

export const UpdateEmailVerificationDtoSchema: z.ZodObject<{
  token: z.ZodString
}> = z.object({
  token: z.string().describe('The password reset token')
})

export type UpdateEmailVerificationDto = z.infer<
  typeof UpdateEmailVerificationDtoSchema
>
