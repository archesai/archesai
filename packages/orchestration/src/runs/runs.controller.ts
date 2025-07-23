import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateRunDtoSchema,
  RUN_ENTITY_KEY,
  RunEntitySchema,
  UpdateRunDtoSchema
} from '@archesai/schemas'

import { createRunRepository } from '#runs/run.repository'
import { createRunsService } from '#runs/runs.service'

export interface RunsPluginOptions {
  databaseService: DatabaseService
  websocketsService: WebsocketsService
}

export const runsController: FastifyPluginAsyncZod<RunsPluginOptions> = async (
  app,
  { databaseService, websocketsService }
) => {
  // Create the run repository and service
  const runRepository = createRunRepository(databaseService)
  const runsService = createRunsService(runRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateRunDtoSchema,
    enableBulkOperations: true,
    entityKey: RUN_ENTITY_KEY,
    entitySchema: RunEntitySchema,
    prefix: '/runs',
    service: runsService,
    updateSchema: UpdateRunDtoSchema
  })
}
