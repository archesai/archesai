import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const UnauthorizedResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z
        .string()
        .describe('You are not authrozied to reach this endpoint.'),
      status: z.string().describe('401'),
      title: z.string().describe('Unauthorized')
    })
  })
  .meta({
    description: 'Schema for 401 Unauthorized response',
    id: 'UnauthorizedResponse'
  })

export type UnauthorizedResponse = z.infer<typeof UnauthorizedResponseSchema>
