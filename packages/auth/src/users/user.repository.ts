import type {
  DrizzleDatabaseService,
  UserSelectModel
} from '@archesai/database'
import type { UserEntity } from '@archesai/schemas'

import { createBaseRepository, UserTable } from '@archesai/database'
import { UserEntitySchema } from '@archesai/schemas'

export const createUserRepository = (
  databaseService: DrizzleDatabaseService
) => {
  return createBaseRepository<UserEntity, UserSelectModel>(
    databaseService,
    UserTable,
    UserEntitySchema
  )
}

export type UserRepository = ReturnType<typeof createUserRepository>
