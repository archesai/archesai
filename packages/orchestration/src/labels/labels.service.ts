import type { WebsocketsService } from '@archesai/core'
import type { LabelEntity } from '@archesai/domain'

import { BaseService } from '@archesai/core'
import { LABEL_ENTITY_KEY } from '@archesai/domain'

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
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname, LABEL_ENTITY_KEY]
    })
  }
}
