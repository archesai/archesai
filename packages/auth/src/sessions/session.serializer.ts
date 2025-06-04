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

  public async deserializeUser(
    id: string,
    done: (err: Error | null, user: UserEntity) => void
  ) {
    const user = await this.usersService.findOne(id)
    done(null, user)
  }

  public serializeUser(
    user: UserEntity,
    done: (err: Error | null, id: string) => void
  ) {
    done(null, user.id)
  }
}
