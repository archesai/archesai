import { Injectable, Logger } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { CreateToolDto } from './dto/create-tool.dto'
import { UpdateToolDto } from './dto/update-tool.dto'
import { ToolEntity, ToolModel } from './entities/tool.entity'
import { ToolRepository } from './tool.repository'

@Injectable()
export class ToolsService extends BaseService<
  ToolEntity,
  CreateToolDto,
  UpdateToolDto,
  ToolRepository,
  ToolModel
> {
  private logger = new Logger(ToolsService.name)
  constructor(
    private toolsRepository: ToolRepository,
    private websocketsService: WebsocketsService
  ) {
    super(toolsRepository)
  }

  async createDefaultTools(orgname: string) {
    await this.toolsRepository.createDefaultTools(orgname)
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit('update', {
      queryKey: ['organizations', orgname, 'tools']
    })
  }

  protected toEntity(model: ToolModel): ToolEntity {
    return new ToolEntity(model)
  }
}
