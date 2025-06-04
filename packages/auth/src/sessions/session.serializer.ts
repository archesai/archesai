import type { ArchesApiRequest } from '@archesai/core'
import type { UserEntity } from '@archesai/domain'

import type { UsersService } from '#users/users.service'

/**
 * Serializer for session management.
 */
export class SessionSerializer {
  private readonly usersService: UsersService

  constructor(usersService: UsersService) {
    this.usersService = usersService
  }

  public async deserializeUser(id: string, _request: ArchesApiRequest) {
    return this.usersService.findOne(id)
  }

  public serializeUser(user: UserEntity, _request: ArchesApiRequest) {
    return Promise.resolve(user.id)
  }
}
