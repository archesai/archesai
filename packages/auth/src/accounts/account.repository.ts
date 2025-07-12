import type { DatabaseService } from '@archesai/core'
import type { AccountInsertModel, AccountSelectModel } from '@archesai/database'
import type { AccountEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { AccountTable } from '@archesai/database'
import { AccountEntitySchema } from '@archesai/schemas'

/**
 * Repository for managing accounts.
 */
export class AccountRepository extends BaseRepository<
  AccountEntity,
  AccountInsertModel,
  AccountSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      AccountEntity,
      AccountInsertModel,
      AccountSelectModel
    >
  ) {
    super(databaseService, AccountTable, AccountEntitySchema)
  }
}
