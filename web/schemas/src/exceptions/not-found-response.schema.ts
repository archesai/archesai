import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const NotFoundResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z.string().describe('The requested resource could not be found.'),
      status: z.string().describe('404'),
      title: z.string().describe('Not Found')
    })
  })
  .meta({
    description: 'Schema for 404 Not Found response',
    id: 'NotFoundResponse'
  })

export type NotFoundResponse = z.infer<typeof NotFoundResponseSchema>
