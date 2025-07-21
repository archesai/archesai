import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const BadRequestResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z.string().describe('The request is invalid or malformed.'),
      status: z.string().describe('400'),
      title: z.string().describe('Bad Request')
    })
  })
  .meta({
    description: 'Schema for 400 Bad Request response',
    id: 'BadRequestResponse'
  })

export type BadRequestResponse = z.infer<typeof BadRequestResponseSchema>
