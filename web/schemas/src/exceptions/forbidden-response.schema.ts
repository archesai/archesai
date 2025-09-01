import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const ForbiddenResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z
        .string()
        .describe('You do not have permission to access this resource.'),
      status: z.string().describe('403'),
      title: z.string().describe('Forbidden')
    })
  })
  .meta({
    description: 'Schema for 403 Forbidden response',
    id: 'ForbiddenResponse'
  })

export type ForbiddenResponse = z.infer<typeof ForbiddenResponseSchema>
