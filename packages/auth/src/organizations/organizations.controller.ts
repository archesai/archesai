import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateOrganizationDtoSchema,
  ORGANIZATION_ENTITY_KEY,
  OrganizationEntitySchema,
  UpdateOrganizationDtoSchema
} from '@archesai/schemas'

import { createOrganizationRepository } from '#organizations/organization.repository'
import { createOrganizationsService } from '#organizations/organizations.service'

export interface OrganizationsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const organizationsPlugin: FastifyPluginAsyncZod<
  OrganizationsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the organization repository and service
  const organizationRepository = createOrganizationRepository(databaseService)
  const organizationsService = createOrganizationsService(
    organizationRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateOrganizationDtoSchema,
    enableBulkOperations: true,
    entityKey: ORGANIZATION_ENTITY_KEY,
    entitySchema: OrganizationEntitySchema,
    prefix: '/organizations',
    service: organizationsService,
    updateSchema: UpdateOrganizationDtoSchema
  })
}
