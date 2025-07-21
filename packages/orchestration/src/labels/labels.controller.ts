import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateLabelDtoSchema,
  LABEL_ENTITY_KEY,
  LabelEntitySchema,
  UpdateLabelDtoSchema
} from '@archesai/schemas'

import { createLabelRepository } from '#labels/label.repository'
import { createLabelsService } from '#labels/labels.service'

export interface LabelsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const labelsController: FastifyPluginAsyncZod<
  LabelsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the label repository and service
  const labelRepository = createLabelRepository(databaseService)
  const labelsService = createLabelsService(labelRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateLabelDtoSchema,
    enableBulkOperations: true,
    entityKey: LABEL_ENTITY_KEY,
    entitySchema: LabelEntitySchema,
    prefix: '/labels',
    service: labelsService,
    updateSchema: UpdateLabelDtoSchema
  })
}
