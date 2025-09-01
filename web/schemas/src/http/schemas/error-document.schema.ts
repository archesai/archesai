import { z } from 'zod'

import { ErrorObjectSchema } from '#http/schemas/error-object.schema'

export const ErrorDocumentSchema: z.ZodObject<{
  error: z.ZodObject<{
    detail: z.ZodString
    status: z.ZodString
    title: z.ZodString
  }>
}> = z.object({
  error: ErrorObjectSchema
})

export type ErrorDocument = z.infer<typeof ErrorDocumentSchema>
