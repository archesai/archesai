import type { WebsocketsService } from '@archesai/core'
import type { SessionEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'

import type { SessionRepository } from '#sessions/session.repository'

export const createSessionsService = (
  sessionRepository: SessionRepository,
  websocketsService: WebsocketsService
) =>
  createBaseService(
    sessionRepository,
    websocketsService,
    emitSessionMutationEvent
  )

const emitSessionMutationEvent = (
  entity: SessionEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.id, 'update', {
    queryKey: ['auth']
  })
}

export type SessionsService = ReturnType<typeof createSessionsService>
