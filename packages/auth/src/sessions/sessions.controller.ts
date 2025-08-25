import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  BetterAuthSessionSchema,
  DocumentSchemaFactory,
  IdParamsSchema,
  SESSION_ENTITY_KEY,
  SessionEntitySchema,
  UnauthorizedResponseSchema
} from '@archesai/schemas'

import type { AuthService } from '#auth/auth.service'

import { createSessionRepository } from '#sessions/session.repository'
import { createSessionsService } from '#sessions/sessions.service'
import { getHeaders } from '#utils/headers'

export interface SessionsControllerOptions {
  authService: AuthService
  databaseService: DatabaseService
  websocketsService: WebsocketsService
}

export const sessionsController: FastifyPluginAsyncZod<
  SessionsControllerOptions
> = async (app, { authService, databaseService, websocketsService }) => {
  // Create the session repository and service
  const sessionRepository = createSessionRepository(databaseService)
  const sessionsService = createSessionsService(
    sessionRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    entityKey: SESSION_ENTITY_KEY,
    entitySchema: SessionEntitySchema,
    prefix: '/sessions',
    service: sessionsService
  })

  app.get(
    `/sessions/current`,
    {
      schema: {
        description: `This endpoint will return the current session information`,
        operationId: 'getSession',
        response: {
          200: BetterAuthSessionSchema,
          401: UnauthorizedResponseSchema
        },
        summary: `Get Session`,
        tags: ['Authentication']
      }
    },
    async (req) => {
      return authService.getSession(getHeaders(req.headers))
    }
  )

  app.patch(
    `/sessions/:id`,
    {
      schema: {
        body: SessionEntitySchema.pick({
          activeOrganizationId: true
        }),
        description: `This endpoint will update the active organization for the current session`,
        operationId: 'updateSession',
        params: IdParamsSchema,
        response: {
          200: DocumentSchemaFactory(SessionEntitySchema),
          401: UnauthorizedResponseSchema
        },
        summary: `Update Session`,
        tags: ['Authentication']
      }
    },
    async (req) => {
      const session = await sessionsService.update(req.params.id, {})
      return {
        data: session
      }
    }
  )
}
