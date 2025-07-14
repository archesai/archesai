import type { WebsocketsService } from '@archesai/core'
import type { UserEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { USER_ENTITY_KEY } from '@archesai/schemas'

import type { UserRepository } from '#users/user.repository'

export const createUsersService = (
  userRepository: UserRepository,
  websocketsService: WebsocketsService
) => createBaseService(userRepository, websocketsService, emitUserMutationEvent)

const emitUserMutationEvent = (
  entity: UserEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.id, 'update', {
    queryKey: ['users', USER_ENTITY_KEY]
  })
}

export type UsersService = ReturnType<typeof createUsersService>
