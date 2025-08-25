import type { BaseService, WebsocketsService } from '@archesai/core'
import type { ToolInsertModel } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { ToolRepository } from '#tools/tool.repository'

export const createToolsService = (
  toolRepository: ToolRepository,
  websocketsService: WebsocketsService
): BaseService<ToolEntity> & {
  createDefaultTools: (organizationId: string) => Promise<{
    count: number
    data: ToolEntity[]
  }>
} => {
  const emitToolMutationEvent = (entity: ToolEntity): void => {
    websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, TOOL_ENTITY_KEY]
    })
  }
  return {
    ...createBaseService(toolRepository, emitToolMutationEvent),
    async createDefaultTools(organizationId: string) {
      const defaultTools: ToolInsertModel[] = [
        {
          description:
            'Extract text from a file. This tool supports all file types.',
          inputMimeType: 'text/plain',
          name: 'Extract Text',
          organizationId,
          outputMimeType: 'text/plain'
        },
        {
          description: 'Create an image from text.',
          inputMimeType: 'text/plain',
          name: 'Text to Image',
          organizationId,
          outputMimeType: 'image/png'
        },
        {
          description: 'Summarize text. This tool supports all languages.',
          inputMimeType: 'text/plain',
          name: 'Summarize',
          organizationId,
          outputMimeType: 'text/plain'
        },
        {
          description:
            'Create embeddings from text. This tool supports all languages.',
          inputMimeType: 'text/plain',
          name: 'Create Embeddings',
          organizationId,
          outputMimeType: 'text/plain' // FIXME make this none
        },
        {
          description:
            'Convert text to speech. This tool supports all languages.',
          inputMimeType: 'text/plain',
          name: 'Text to Speech',
          organizationId,
          outputMimeType: 'audio/mpeg'
        }
      ]

      const result = await toolRepository.createMany(defaultTools)
      result.data.forEach((tool) => {
        emitToolMutationEvent(tool)
      })
      return result
    }
  }
}

export type ToolsService = ReturnType<typeof createToolsService>
