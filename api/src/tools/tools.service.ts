import { Injectable } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { ToolEntity, ToolModel } from './entities/tool.entity'
import { ToolRepository } from './tool.repository'

@Injectable()
export class ToolsService extends BaseService<
  ToolEntity,
  ToolModel,
  ToolRepository
> {
  constructor(
    private toolsRepository: ToolRepository,
    private websocketsService: WebsocketsService
  ) {
    super(toolsRepository)
  }

  async createDefaultTools(orgname: string) {
    await this.toolsRepository.createDefaultTools(orgname)
  }

  protected emitMutationEvent(entity: ToolEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'tools']
    })
  }

  protected toEntity(model: ToolModel): ToolEntity {
    return new ToolEntity(model)
  }
}
