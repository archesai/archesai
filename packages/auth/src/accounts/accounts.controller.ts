import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import type { WebsocketsService } from '@archesai/core'
import type { DrizzleDatabaseService } from '@archesai/database'

import { crudPlugin } from '@archesai/core'
import { ACCOUNT_ENTITY_KEY, AccountEntitySchema } from '@archesai/schemas'

import { createAccountRepository } from '#accounts/account.repository'
import { createAccountsService } from '#accounts/accounts.service'

export interface AccountsPluginOptions {
  databaseService: DrizzleDatabaseService
  websocketsService: WebsocketsService
}

export const accountsPlugin: FastifyPluginAsyncZod<
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
    entityKey: ACCOUNT_ENTITY_KEY,
    entitySchema: AccountEntitySchema,
    prefix: '/accounts',
    service: accountsService
  })
}
