import { z } from 'zod'

export const CreateEmailVerificationDtoSchema: z.ZodObject<{
  email: z.ZodEmail
  userId: z.ZodUUID
}> = z.object({
  email: z.email().describe('The e-mail to send the confirmation token to'),
  userId: z
    .uuid()
    .describe('The user ID of the user requesting the email verification')
})

export type CreateEmailVerificationDto = z.infer<
  typeof CreateEmailVerificationDtoSchema
>
