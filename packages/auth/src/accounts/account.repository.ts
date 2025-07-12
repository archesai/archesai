import type { DatabaseService } from '@archesai/core'
import type { AccountEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { AccountTable } from '@archesai/database'
import { AccountEntitySchema } from '@archesai/schemas'

/**
 * Repository for managing accounts.
 */
export class AccountRepository extends BaseRepository<AccountEntity> {
  constructor(databaseService: DatabaseService<AccountEntity>) {
    super(databaseService, AccountTable, AccountEntitySchema)
  }
}
