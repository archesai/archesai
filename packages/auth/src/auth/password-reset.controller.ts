import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import {
  CreatePasswordResetDtoSchema,
  NoContentResponseSchema,
  NotFoundResponseSchema,
  UnauthorizedResponseSchema,
  UpdatePasswordResetDtoSchema
} from '@archesai/schemas'

export const passwordResetController: FastifyPluginAsyncZod = async (app) => {
  app.post(
    `/password-reset/verify`,
    {
      schema: {
        body: UpdatePasswordResetDtoSchema,
        description:
          'This endpoint will verify your password change with a token',
        operationId: 'confirmPasswordReset',
        response: {
          204: NoContentResponseSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Verify password reset',
        tags: ['Password Reset']
      }
    },
    () => {
      throw new Error(
        'Password reset is not implemented yet. Please use the password reset request endpoint.'
      )
    }
  )

  app.post(
    `/password-reset/request`,
    {
      schema: {
        body: CreatePasswordResetDtoSchema,
        description: 'This endpoint will request a password reset link',
        operationId: 'requestPasswordReset',
        response: {
          204: NoContentResponseSchema
        },
        summary: 'Request password reset',
        tags: ['Password Reset']
      }
    },
    () => {
      throw new Error(
        'Password reset is not implemented yet. Please use the password reset request endpoint.'
      )
    }
  )

  await Promise.resolve()
}
