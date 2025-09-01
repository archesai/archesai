import { z } from 'zod'

export const ErrorObjectSchema: z.ZodObject<{
  detail: z.ZodString
  status: z.ZodString
  title: z.ZodString
}> = z.object({
  detail: z.string().describe('The requested resource does not exist.'),
  status: z.string().describe('404'),
  title: z.string().describe('Not Found')
})

export type ErrorObject = z.infer<typeof ErrorObjectSchema>
