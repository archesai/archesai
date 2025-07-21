import { z } from 'zod'

import type { ErrorDocument } from '#http/schemas/error-document.schema'

export const ValidationErrorResponseSchema: z.ZodType<ErrorDocument> = z
  .object({
    error: z.object({
      detail: z.string().describe('Validation failed for one or more fields.'),
      details: z.array(
        z.object({
          field: z.string().describe('username, email'),
          message: z
            .string()
            .describe('Username is required., Email format is invalid.'),
          value: z.string().optional().describe('john_doe, invalid-email')
        })
      ),
      status: z.string().describe('422'),
      title: z.string().describe('Validation Error')
    })
  })
  .meta({
    description: 'Schema for 422 Validation Error response',
    id: 'ValidationErrorResponse'
  })

export type ValidationErrorResponse = z.infer<
  typeof ValidationErrorResponseSchema
>
