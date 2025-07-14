import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  ACCOUNT_ENTITY_KEY,
  AccountEntitySchema,
  CreateAccountDtoSchema
} from '@archesai/schemas'

import { createAccountRepository } from '#accounts/account.repository'
import { createAccountsService } from '#accounts/accounts.service'

export interface AccountsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const accountsPlugin: FastifyPluginAsyncTypebox<
  AccountsPluginOptions
> = async (app, { databaseService, websocketsService }) => {
  // Create the account repository and service
  const accountRepository = createAccountRepository(databaseService)
  const accountsService = createAccountsService(
    accountRepository,
    websocketsService
  )

  // Register CRUD routes
  await app.register(crudPlugin, {
    createSchema: CreateAccountDtoSchema,
    enableBulkOperations: true,
    entityKey: ACCOUNT_ENTITY_KEY,
    entitySchema: AccountEntitySchema,
    prefix: '/accounts',
    service: accountsService,
    updateSchema: AccountEntitySchema
  })
}
