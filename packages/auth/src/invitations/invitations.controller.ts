import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type {
  InvitationInsertModel,
  InvitationSelectModel
} from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateInvitationDtoSchema,
  InvitationEntitySchema,
  TOOL_ENTITY_KEY,
  UpdateInvitationDtoSchema
} from '@archesai/schemas'

import { createInvitationRepository } from '#invitations/invitation.repository'
import { createInvitationsService } from '#invitations/invitations.service'

export interface InvitationsPluginOptions {
  databaseService: DatabaseService<InvitationInsertModel, InvitationSelectModel>
  websocketsService: WebsocketsService
}

export const invitationsPlugin: FastifyPluginAsyncTypebox<
  InvitationsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the invitation repository and service
  const invitationRepository = createInvitationRepository(databaseService)
  const invitationsService = createInvitationsService(
    invitationRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateInvitationDtoSchema,
    enableBulkOperations: true,
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: InvitationEntitySchema,
    prefix: '/invitations',
    service: invitationsService,
    updateSchema: UpdateInvitationDtoSchema
  })
}
