import type { ModuleMetadata } from '@archesai/core'
import type { AccountInsertModel, AccountSelectModel } from '@archesai/database'
import type { AccountEntity } from '@archesai/schemas'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { AccountRepository } from '#accounts/account.repository'
import { AccountsService } from '#accounts/accounts.service'
import { HashingModule } from '#hashing/hashing.module'
import { HashingService } from '#hashing/hashing.service'

export const AccountsModuleDefinition: ModuleMetadata = {
  exports: [AccountsService],
  imports: [DatabaseModule, HashingModule, WebsocketsModule],
  providers: [
    {
      inject: [AccountRepository, HashingService, WebsocketsService],
      provide: AccountsService,
      useFactory: (
        accountRepository: AccountRepository,
        hashingService: HashingService,
        websocketsService: WebsocketsService
      ) =>
        new AccountsService(
          accountRepository,
          hashingService,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: AccountRepository,
      useFactory: (
        databaseService: DatabaseService<
          AccountEntity,
          AccountInsertModel,
          AccountSelectModel
        >
      ) => new AccountRepository(databaseService)
    }
  ]
}

export const AccountsModule = (() =>
  createModule(class AccountsModule {}, AccountsModuleDefinition))()
