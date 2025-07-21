import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  ARTIFACT_ENTITY_KEY,
  ArtifactEntitySchema,
  CreateArtifactDtoSchema,
  UpdateArtifactDtoSchema
} from '@archesai/schemas'

import { createArtifactRepository } from '#artifacts/artifact.repository'
import { createArtifactsService } from '#artifacts/artifacts.service'

export interface ArtifactsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const artifactsController: FastifyPluginAsyncZod<
  ArtifactsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  app.log.info('Registering artifacts controller')
  // Create the artifact repository and service
  const artifactRepository = createArtifactRepository(databaseService)
  const artifactsService = createArtifactsService(
    artifactRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateArtifactDtoSchema,
    enableBulkOperations: true,
    entityKey: ARTIFACT_ENTITY_KEY,
    entitySchema: ArtifactEntitySchema,
    prefix: '/artifacts',
    service: artifactsService,
    updateSchema: UpdateArtifactDtoSchema
  })
}
