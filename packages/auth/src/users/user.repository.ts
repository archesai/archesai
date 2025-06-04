import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { UserTable } from '@archesai/database'
import { UserEntity } from '@archesai/domain'

/**
 * Repository for handling users.
 */
export class UserRepository extends BaseRepository<UserEntity> {
  constructor(databaseService: DatabaseService<UserEntity>) {
    super(databaseService, UserTable, UserEntity)
  }
}
