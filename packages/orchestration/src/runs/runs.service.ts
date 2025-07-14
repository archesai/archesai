import type { WebsocketsService } from '@archesai/core'
import type { RunEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { RunRepository } from '#runs/run.repository'

export const createRunsService = (
  runRepository: RunRepository,
  websocketsService: WebsocketsService
) => createBaseService(runRepository, websocketsService, emitRunMutationEvent)

const emitRunMutationEvent = (
  entity: RunEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, TOOL_ENTITY_KEY]
  })
}

export type RunsService = ReturnType<typeof createRunsService>
