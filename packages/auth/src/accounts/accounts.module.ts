import type { ModuleMetadata } from '@archesai/core'
import type { AccountEntity } from '@archesai/domain'

import {
  DatabaseModule,
  DatabaseService,
  Module,
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
      useFactory: (databaseService: DatabaseService<AccountEntity>) =>
        new AccountRepository(databaseService)
    }
  ]
}

@Module(AccountsModuleDefinition)
export class AccountsModule {}
