import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const InternalServerErrorResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z
        .string()
        .describe('An unexpected error occurred on the server.'),
      status: z.string().describe('500'),
      title: z.string().describe('Internal Server Error')
    })
  })
  .meta({
    description: 'Schema for 500 Internal Server Error response',
    id: 'InternalServerErrorResponse'
  })

export type InternalServerErrorResponse = z.infer<
  typeof InternalServerErrorResponseSchema
>
