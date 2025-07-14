import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  CreateUserDtoSchema,
  UpdateUserDtoSchema,
  USER_ENTITY_KEY,
  UserEntitySchema
} from '@archesai/schemas'

import { createUserRepository } from '#users/user.repository'
import { createUsersService } from '#users/users.service'

export interface UsersPluginOptions {
  databaseService: DrizzleDatabaseService
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
    entityKey: USER_ENTITY_KEY,
    entitySchema: UserEntitySchema,
    prefix: '/users',
    service: usersService,
    updateSchema: UpdateUserDtoSchema
  })
}
