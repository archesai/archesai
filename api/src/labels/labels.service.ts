import { Injectable } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { LabelEntity, LabelModel } from './entities/label.entity'
import { LabelRepository } from './label.repository'

@Injectable()
export class LabelsService extends BaseService<
  LabelEntity,
  LabelModel,
  LabelRepository
> {
  constructor(
    private labelRepository: LabelRepository,
    private websocketsService: WebsocketsService
  ) {
    super(labelRepository)
  }

  protected emitMutationEvent(entity: LabelEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'labels']
    })
  }

  protected toEntity(model: LabelModel): LabelEntity {
    return new LabelEntity(model)
  }
}
