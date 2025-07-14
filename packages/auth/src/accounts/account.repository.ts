import type {
  AccountSelectModel,
  DrizzleDatabaseService
} from '@archesai/database'
import type { AccountEntity } from '@archesai/schemas'

import { AccountTable, createBaseRepository } from '@archesai/database'
import { AccountEntitySchema } from '@archesai/schemas'

export const createAccountRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<AccountEntity, AccountSelectModel>(
    databaseService,
    AccountTable,
    AccountEntitySchema
  )
}

export type AccountRepository = ReturnType<typeof createAccountRepository>
