import { z } from 'zod'

export const UpdateEmailChangeDtoSchema: z.ZodObject<{
  newEmail: z.ZodEmail
  token: z.ZodString
  userId: z.ZodUUID
}> = z.object({
  newEmail: z.email().describe('The e-mail to send the confirmation token to'),
  token: z.string().describe('The password reset token'),
  userId: z
    .uuid()
    .describe('The user ID of the user requesting the email change')
})

export type UpdateEmailChangeDto = z.infer<typeof UpdateEmailChangeDtoSchema>
