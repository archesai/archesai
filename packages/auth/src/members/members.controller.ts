import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateMemberDtoSchema,
  MEMBER_ENTITY_KEY,
  MemberEntitySchema,
  UpdateMemberDtoSchema
} from '@archesai/schemas'

import { createMemberRepository } from '#members/member.repository'
import { createMembersService } from '#members/members.service'

export interface MembersPluginOptions {
  databaseService: DatabaseService
  websocketsService: WebsocketsService
}

export const membersPlugin: FastifyPluginAsyncZod<
  MembersPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the member repository and service
  const memberRepository = createMemberRepository(databaseService)
  const membersService = createMembersService(
    memberRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateMemberDtoSchema,
    enableBulkOperations: true,
    entityKey: MEMBER_ENTITY_KEY,
    entitySchema: MemberEntitySchema,
    prefix: '/members',
    service: membersService,
    updateSchema: UpdateMemberDtoSchema
  })
}
