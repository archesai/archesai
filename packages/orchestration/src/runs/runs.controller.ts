import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

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
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const runsController: FastifyPluginAsyncTypebox<
  RunsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
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
