import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import {
  BetterAuthSessionSchema,
  NoContentResponseSchema,
  NotFoundResponseSchema,
  UnauthorizedResponseSchema,
  UpdateEmailVerificationDtoSchema
} from '@archesai/schemas'

export const emailVerificationController: FastifyPluginAsyncZod = async (
  app
) => {
  app.post(
    `/email-verification/verify`,
    {
      schema: {
        body: UpdateEmailVerificationDtoSchema,
        description: 'This endpoint will confirm your e-mail with a token',
        operationId: 'confirmEmailVerification',
        response: {
          200: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Confirm e-mail verification',
        tags: ['Email Verification']
      }
    },
    () => {
      throw new Error(
        'Email verification is not implemented yet. Please use the email verification request endpoint.'
      )
    }
  )

  app.post(
    `/email-verification/request`,
    {
      schema: {
        description:
          'This endpoint will send an e-mail verification link to you. ADMIN ONLY.',
        operationId: 'requestEmailVerification',
        response: {
          204: NoContentResponseSchema
        },
        security: [{ bearerAuth: [] }], // ✅ add this line
        summary: 'Request e-mail verification',
        tags: ['Email Verification']
      }
    },
    () => {
      throw new Error(
        'Email verification is not implemented yet. Please use the email verification request endpoint.'
      )
    }
  )

  await Promise.resolve()
}
