import type { UserEntity } from '#auth/users/entities/user.entity'

export class UserCreatedEvent {
  public user: UserEntity

  constructor(user: UserEntity) {
    this.user = user
  }
}
