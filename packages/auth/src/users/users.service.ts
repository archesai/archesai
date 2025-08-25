import type { BaseService, WebsocketsService } from '@archesai/core'
import type { UserEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { USER_ENTITY_KEY } from '@archesai/schemas'

import type { UserRepository } from '#users/user.repository'

export const createUsersService = (
  userRepository: UserRepository,
  websocketsService: WebsocketsService
): BaseService<UserEntity> => {
  const emitUserMutationEvent = (entity: UserEntity): void => {
    websocketsService.broadcastEvent(entity.id, 'update', {
      queryKey: ['users', USER_ENTITY_KEY]
    })
  }
  return createBaseService(userRepository, emitUserMutationEvent)
}

export type UsersService = ReturnType<typeof createUsersService>
