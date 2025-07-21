import { z } from 'zod'

export const CreateAccountDtoSchema: z.ZodObject<{
  email: z.ZodEmail
  name: z.ZodString
  password: z.ZodString
}> = z.object({
  email: z.email().describe('The email address associated with the account'),
  name: z.string().min(1).describe('The name of the user creating the account'),
  password: z.string().describe('The password for the account')
})

export type CreateAccountDto = z.infer<typeof CreateAccountDtoSchema>
