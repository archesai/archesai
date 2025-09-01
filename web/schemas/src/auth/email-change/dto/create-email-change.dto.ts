import { z } from 'zod'

export const CreateEmailChangeDtoSchema: z.ZodObject<{
  newEmail: z.ZodEmail
  userId: z.ZodUUID
}> = z.object({
  newEmail: z.email().describe('The e-mail to send the confirmation token to'),
  userId: z
    .uuid()
    .describe('The user ID of the user requesting the email change')
})

export type CreateEmailChangeDto = z.infer<typeof CreateEmailChangeDtoSchema>
