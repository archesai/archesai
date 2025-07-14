import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type {
  OrganizationInsertModel,
  OrganizationSelectModel
} from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateOrganizationDtoSchema,
  OrganizationEntitySchema,
  TOOL_ENTITY_KEY,
  UpdateOrganizationDtoSchema
} from '@archesai/schemas'

import { createOrganizationRepository } from '#organizations/organization.repository'
import { createOrganizationsService } from '#organizations/organizations.service'

export interface OrganizationsPluginOptions {
  databaseService: DatabaseService<
    OrganizationInsertModel,
    OrganizationSelectModel
  >
  websocketsService: WebsocketsService
}

export const organizationsPlugin: FastifyPluginAsyncTypebox<
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
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: OrganizationEntitySchema,
    prefix: '/organizations',
    service: organizationsService,
    updateSchema: UpdateOrganizationDtoSchema
  })
}
