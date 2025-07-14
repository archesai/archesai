import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type { MemberInsertModel, MemberSelectModel } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateMemberDtoSchema,
  MemberEntitySchema,
  TOOL_ENTITY_KEY,
  UpdateMemberDtoSchema
} from '@archesai/schemas'

import { createMemberRepository } from '#members/member.repository'
import { createMembersService } from '#members/members.service'

export interface MembersPluginOptions {
  databaseService: DatabaseService<MemberInsertModel, MemberSelectModel>
  websocketsService: WebsocketsService
}

export const membersPlugin: FastifyPluginAsyncTypebox<
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
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: MemberEntitySchema,
    prefix: '/members',
    service: membersService,
    updateSchema: UpdateMemberDtoSchema
  })
}
