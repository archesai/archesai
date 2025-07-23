import type { WebsocketsService } from '@archesai/core'
import type { SessionEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'

import type { SessionRepository } from '#sessions/session.repository'

export const createSessionsService = (
  sessionRepository: SessionRepository,
  websocketsService: WebsocketsService
) => {
  const emitSessionMutationEvent = (entity: SessionEntity): void => {
    websocketsService.broadcastEvent(entity.id, 'update', {
      queryKey: ['auth']
    })
  }
  return createBaseService(sessionRepository, emitSessionMutationEvent)
}

export type SessionsService = ReturnType<typeof createSessionsService>
