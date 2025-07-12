import type { DatabaseService } from '@archesai/core'
import type { UserEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { UserTable } from '@archesai/database'
import { UserEntitySchema } from '@archesai/schemas'

/**
 * Repository for handling users.
 */
export class UserRepository extends BaseRepository<UserEntity> {
  constructor(databaseService: DatabaseService<UserEntity>) {
    super(databaseService, UserTable, UserEntitySchema)
  }
}
