import type { WebsocketsService } from '@archesai/core'
import type { LabelEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { LABEL_ENTITY_KEY } from '@archesai/schemas'

import type { LabelRepository } from '#labels/label.repository'

/**
 * Service for labels.
 */
export class LabelsService extends BaseService<LabelEntity> {
  private readonly websocketsService: WebsocketsService

  constructor(
    labelRepository: LabelRepository,
    websocketsService: WebsocketsService
  ) {
    super(labelRepository)
    this.websocketsService = websocketsService
  }

  protected emitMutationEvent(entity: LabelEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, LABEL_ENTITY_KEY]
    })
  }
}
