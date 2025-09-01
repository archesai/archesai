import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const ConflictResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z
        .string()
        .describe(
          'The request conflicts with the current state of the resource.'
        ),
      status: z.string().describe('409'),
      title: z.string().describe('Conflict')
    })
  })
  .meta({
    description: 'Schema for 409 Conflict response',
    id: 'ConflictResponse'
  })

export type ConflictResponse = z.infer<typeof ConflictResponseSchema>
