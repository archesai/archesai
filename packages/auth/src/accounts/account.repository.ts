import type { DatabaseService } from '@archesai/core'
import type { AccountInsertModel, AccountSelectModel } from '@archesai/database'
import type { AccountEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { AccountTable } from '@archesai/database'
import { AccountEntitySchema } from '@archesai/schemas'

export const createAccountRepository = (
  databaseService: DatabaseService<AccountInsertModel, AccountSelectModel>
) => {
  return createBaseRepository<
    AccountEntity,
    AccountInsertModel,
    AccountSelectModel
  >(databaseService, AccountTable, AccountEntitySchema)
}

export type AccountRepository = ReturnType<typeof createAccountRepository>
