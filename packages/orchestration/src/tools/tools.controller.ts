import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateToolDtoSchema,
  TOOL_ENTITY_KEY,
  ToolEntitySchema,
  UpdateToolDtoSchema
} from '@archesai/schemas'

import { createToolRepository } from '#tools/tool.repository'
import { createToolsService } from '#tools/tools.service'

export interface ToolsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const toolsController: FastifyPluginAsyncTypebox<
  ToolsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the tool repository and service
  const toolRepository = createToolRepository(databaseService)
  const toolsService = createToolsService(toolRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateToolDtoSchema,
    enableBulkOperations: true,
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: ToolEntitySchema,
    prefix: '/tools',
    service: toolsService,
    updateSchema: UpdateToolDtoSchema
  })
}
