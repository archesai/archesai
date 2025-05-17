import type { EventBus, WebsocketsService } from '@archesai/core'
import type { BaseInsertion, UserEntity } from '@archesai/domain'

import { BaseService } from '@archesai/core'
import { USER_ENTITY_KEY, UserCreatedEvent } from '@archesai/domain'

import type { UserRepository } from '#users/user.repository'

/**
 * Service for handling users.
 */
export class UsersService extends BaseService<UserEntity> {
  private readonly eventBus: EventBus
  private readonly userRepository: UserRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    eventBus: EventBus,
    userRepository: UserRepository,
    websocketsService: WebsocketsService
  ) {
    super(userRepository)
    this.eventBus = eventBus
    this.userRepository = userRepository
    this.websocketsService = websocketsService
  }

  public override async create(value: BaseInsertion<UserEntity>) {
    const user = await this.userRepository.create({
      ...value,
      orgname:
        value.orgname ||
        (value.email.split('@')[0] ?? value.email) +
          '-' +
          Math.random().toString(36).substring(2, 6)
    })
    this.eventBus.emit('user.created', new UserCreatedEvent(user))
    return user
  }

  public async deactivate(id: string): Promise<void> {
    await this.userRepository.deactivate(id)
  }

  public async findOneByEmail(email: string): Promise<UserEntity> {
    return this.userRepository.findOneByEmail(email)
  }

  protected emitMutationEvent(user: UserEntity): void {
    this.websocketsService.broadcastEvent(user.orgname, 'update', {
      queryKey: [USER_ENTITY_KEY]
    })
  }
}
