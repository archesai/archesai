import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
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

export const usersPlugin: FastifyPluginAsyncZod<UsersPluginOptions> = async (
  app,
  { databaseService, websocketsService }
) => {
  // Create the user repository and service
  const userRepository = createUserRepository(databaseService)
  const usersService = createUsersService(userRepository, websocketsService)

  // Register CRUD routes
  await app.register(crudPlugin, {
    enableBulkOperations: true,
    entityKey: USER_ENTITY_KEY,
    entitySchema: UserEntitySchema,
    prefix: '/users',
    service: usersService,
    updateSchema: UpdateUserDtoSchema
  })
}
