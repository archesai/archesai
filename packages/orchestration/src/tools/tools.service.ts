import type { WebsocketsService } from '@archesai/core'
import type { ToolEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { ToolRepository } from '#tools/tool.repository'

/**
 * Service for tools.
 */
export class ToolsService extends BaseService<ToolEntity> {
  private readonly toolsRepository: ToolRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    toolsRepository: ToolRepository,
    websocketsService: WebsocketsService
  ) {
    super(toolsRepository)
    this.toolsRepository = toolsRepository
    this.websocketsService = websocketsService
  }

  public async createDefaultTools(orgname: string) {
    return this.toolsRepository.createDefaultTools(orgname)
  }

  protected emitMutationEvent(entity: ToolEntity): void {
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname, TOOL_ENTITY_KEY]
    })
  }
}
