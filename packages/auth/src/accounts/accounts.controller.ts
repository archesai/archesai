import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { DatabaseService, WebsocketsService } from '@archesai/core'
import type { AccountInsertModel, AccountSelectModel } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import {
  AccountEntitySchema,
  CreateAccountDtoSchema,
  TOOL_ENTITY_KEY
} from '@archesai/schemas'

import { createAccountRepository } from '#accounts/account.repository'
import { createAccountsService } from '#accounts/accounts.service'

export interface AccountsPluginOptions {
  databaseService: DatabaseService<AccountInsertModel, AccountSelectModel>
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
    entityKey: TOOL_ENTITY_KEY,
    entitySchema: AccountEntitySchema,
    prefix: '/accounts',
    service: accountsService,
    updateSchema: AccountEntitySchema
  })
}
