import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateInvitationDtoSchema,
  INVITATION_ENTITY_KEY,
  InvitationEntitySchema,
  UpdateInvitationDtoSchema
} from '@archesai/schemas'

import { createInvitationRepository } from '#invitations/invitation.repository'
import { createInvitationsService } from '#invitations/invitations.service'

export interface InvitationsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const invitationsPlugin: FastifyPluginAsyncZod<
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
    entityKey: INVITATION_ENTITY_KEY,
    entitySchema: InvitationEntitySchema,
    prefix: '/invitations',
    service: invitationsService,
    updateSchema: UpdateInvitationDtoSchema
  })
}
