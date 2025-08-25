import type { BaseService, WebsocketsService } from '@archesai/core'
import type { LabelEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { LabelRepository } from '#labels/label.repository'

export const createLabelsService = (
  labelRepository: LabelRepository,
  websocketsService: WebsocketsService
): BaseService<LabelEntity> => {
  const emitLabelMutationEvent = (entity: LabelEntity): void => {
    websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, TOOL_ENTITY_KEY]
    })
  }
  return createBaseService(labelRepository, emitLabelMutationEvent)
}

export type LabelsService = ReturnType<typeof createLabelsService>
