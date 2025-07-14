import type { WebsocketsService } from '@archesai/core'
import type { PipelineEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { PipelineRepository } from '#pipelines/pipeline.repository'

export const createPipelinesService = (
  pipelineRepository: PipelineRepository,
  websocketsService: WebsocketsService
) =>
  createBaseService(
    pipelineRepository,
    websocketsService,
    emitPipelineMutationEvent
  )

const emitPipelineMutationEvent = (
  entity: PipelineEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, TOOL_ENTITY_KEY]
  })
}

export type PipelinesService = ReturnType<typeof createPipelinesService>
