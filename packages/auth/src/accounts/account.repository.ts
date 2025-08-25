import type { BaseRepository } from '@archesai/core'
import type { DatabaseService } from '@archesai/database'
import type { AccountEntity } from '@archesai/schemas'

import { AccountTable, createBaseRepository } from '@archesai/database'
import { AccountEntitySchema } from '@archesai/schemas'

export const createAccountRepository = (
  databaseService: DatabaseService
): BaseRepository<
  AccountEntity,
  (typeof AccountTable)['$inferInsert'],
  (typeof AccountTable)['$inferSelect']
> => {
  return createBaseRepository<AccountEntity>(
    databaseService,
    AccountTable,
    AccountEntitySchema
  )
}

export type AccountRepository = ReturnType<typeof createAccountRepository>
