import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type { UserInsertModel, UserSelectModel } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateUserDtoSchema,
  TOOL_ENTITY_KEY,
  UpdateUserDtoSchema,
  UserEntitySchema
} from '@archesai/schemas'

import { createUserRepository } from '#users/user.repository'
import { createUsersService } from '#users/users.service'

export interface UsersPluginOptions {
  databaseService: DatabaseService<UserInsertModel, UserSelectModel>
  websocketsService: WebsocketsService
}

export const usersPlugin: FastifyPluginAsyncTypebox<
  UsersPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the user repository and service
  const userRepository = createUserRepository(databaseService)
  const usersService = createUsersService(userRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateUserDtoSchema,
    enableBulkOperations: true,
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: UserEntitySchema,
    prefix: '/users',
    service: usersService,
    updateSchema: UpdateUserDtoSchema
  })
}
