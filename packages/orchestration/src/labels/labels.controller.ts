import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type { LabelInsertModel, LabelSelectModel } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateLabelDtoSchema,
  LabelEntitySchema,
  TOOL_ENTITY_KEY,
  UpdateLabelDtoSchema
} from '@archesai/schemas'

import { createLabelRepository } from '#labels/label.repository'
import { createLabelsService } from '#labels/labels.service'

export interface LabelsPluginOptions {
  databaseService: DatabaseService<LabelInsertModel, LabelSelectModel>
  websocketsService: WebsocketsService
}

export const labelsController: FastifyPluginAsyncTypebox<
  LabelsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the label repository and service
  const labelRepository = createLabelRepository(databaseService)
  const labelsService = createLabelsService(labelRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateLabelDtoSchema,
    enableBulkOperations: true,
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: LabelEntitySchema,
    prefix: '/labels',
    service: labelsService,
    updateSchema: UpdateLabelDtoSchema
  })
}
