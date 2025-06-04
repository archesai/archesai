import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { AccountTable } from '@archesai/database'
import { AccountEntity } from '@archesai/domain'

/**
 * Repository for managing accounts.
 */
export class AccountRepository extends BaseRepository<AccountEntity> {
  constructor(databaseService: DatabaseService<AccountEntity>) {
    super(databaseService, AccountTable, AccountEntity)
  }
}
