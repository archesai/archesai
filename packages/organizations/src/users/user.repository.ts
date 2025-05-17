import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { USER_ENTITY_KEY, UserEntity } from '@archesai/domain'

/**
 * Repository for handling users.
 */
export class UserRepository extends BaseRepository<UserEntity> {
  constructor(databaseService: DatabaseService<UserEntity>) {
    super(databaseService, USER_ENTITY_KEY, UserEntity)
  }

  public async deactivate(id: string): Promise<void> {
    await this.update(id, {
      deactivated: true
    })
  }

  public async findOneByEmail(email: string): Promise<UserEntity> {
    return this.findFirst({
      filter: {
        email: {
          equals: email
        }
      },
      page: {
        number: 1,
        size: 1
      },
      sort: '-createdAt'
    })
  }
}
