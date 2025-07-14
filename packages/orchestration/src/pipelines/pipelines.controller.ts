import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreatePipelineDtoSchema,
  PIPELINE_ENTITY_KEY,
  PipelineEntitySchema,
  UpdatePipelineDtoSchema
} from '@archesai/schemas'

import { createPipelineRepository } from '#pipelines/pipeline.repository'
import { createPipelinesService } from '#pipelines/pipelines.service'

export interface PipelinesPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const pipelinesController: FastifyPluginAsyncTypebox<
  PipelinesPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the pipeline repository and service
  const pipelineRepository = createPipelineRepository(databaseService)
  const pipelinesService = createPipelinesService(
    pipelineRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreatePipelineDtoSchema,
    enableBulkOperations: true,
    entityKey: PIPELINE_ENTITY_KEY,
    entitySchema: PipelineEntitySchema,
    prefix: '/pipelines',
    service: pipelinesService,
    updateSchema: UpdatePipelineDtoSchema
  })
}
