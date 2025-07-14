import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type { SessionInsertModel, SessionSelectModel } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import { SessionEntitySchema, TOOL_ENTITY_KEY } from '@archesai/schemas'

import { createSessionRepository } from '#sessions/session.repository'
import { createSessionsService } from '#sessions/sessions.service'

export interface SessionsPluginOptions {
  databaseService: DatabaseService<SessionInsertModel, SessionSelectModel>
  websocketsService: WebsocketsService
}

export const sessionsPlugin: FastifyPluginAsyncTypebox<
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
    createSchema: SessionEntitySchema,
    enableBulkOperations: true,
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: SessionEntitySchema,
    prefix: '/sessions',
    service: sessionsService,
    updateSchema: SessionEntitySchema
  })
}
