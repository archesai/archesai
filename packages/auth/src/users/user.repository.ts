import type { DatabaseService } from '@archesai/core'
import type { UserInsertModel, UserSelectModel } from '@archesai/database'
import type { UserEntity } from '@archesai/schemas'

import { createBaseRepository } from '@archesai/core'
import { UserTable } from '@archesai/database'
import { UserEntitySchema } from '@archesai/schemas'

export const createUserRepository = (
  databaseService: DatabaseService<UserInsertModel, UserSelectModel>
) => {
  return createBaseRepository<UserEntity, UserInsertModel, UserSelectModel>(
    databaseService,
    UserTable,
    UserEntitySchema
  )
}

export type UserRepository = ReturnType<typeof createUserRepository>
