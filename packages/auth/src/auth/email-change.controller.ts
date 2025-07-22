import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import {
  CreateEmailChangeDtoSchema,
  NoContentResponseSchema,
  NotFoundResponseSchema,
  UnauthorizedResponseSchema,
  UpdateEmailChangeDtoSchema
} from '@archesai/schemas'

export const emailChangeController: FastifyPluginAsyncZod = async (app) => {
  app.post(
    `/email-change/verify`,
    {
      schema: {
        body: UpdateEmailChangeDtoSchema,
        description:
          'This endpoint will verify your e-mail change with a token',
        operationId: 'confirmEmailChange',
        response: {
          204: NoContentResponseSchema,
          401: UnauthorizedResponseSchema,
          404: NotFoundResponseSchema
        },
        summary: 'Verify e-mail change',
        tags: ['Email Change']
      }
    },
    () => {
      throw new Error(
        'Email change is not implemented yet. Please use the email change request endpoint.'
      )
    }
  )

  app.post(
    `/email-change/request`,
    {
      schema: {
        body: CreateEmailChangeDtoSchema,
        description:
          'This endpoint will request your e-mail change with a token',
        operationId: 'requestEmailChange',
        response: {
          204: NoContentResponseSchema
        },
        summary: 'Request e-mail change',
        tags: ['Email Change']
      }
    },
    () => {
      throw new Error(
        'Email change is not implemented yet. Please use the email change request endpoint.'
      )
    }
  )

  await Promise.resolve()
}
