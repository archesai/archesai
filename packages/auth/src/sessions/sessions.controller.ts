import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import { SESSION_ENTITY_KEY, SessionEntitySchema } from '@archesai/schemas'

import { createSessionRepository } from '#sessions/session.repository'
import { createSessionsService } from '#sessions/sessions.service'

export interface SessionsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const sessionsPlugin: FastifyPluginAsyncZod<
  SessionsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
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
}
