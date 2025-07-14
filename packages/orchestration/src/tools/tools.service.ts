import type { WebsocketsService } from '@archesai/core'
import type { ToolInsertModel } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { TOOL_ENTITY_KEY } from '@archesai/schemas'

import type { ToolRepository } from '#tools/tool.repository'

export const createToolsService = (
  toolRepository: ToolRepository,
  websocketsService: WebsocketsService
) => {
  return {
    ...createBaseService(
      toolRepository,
      websocketsService,
      emitToolMutationEvent
    ),
    async createDefaultTools(organizationId: string) {
      const defaultTools: ToolInsertModel[] = [
        {
          description:
            'Extract text from a file. This tool supports all file types.',
          inputType: 'TEXT',
          name: 'Extract Text',
          organizationId,
          outputType: 'TEXT',
          toolBase: 'extract-text'
        },
        {
          description: 'Create an image from text.',
          inputType: 'TEXT',
          name: 'Text to Image',
          organizationId,
          outputType: 'IMAGE',
          toolBase: 'text-to-image'
        },
        {
          description: 'Summarize text. This tool supports all languages.',
          inputType: 'TEXT',
          name: 'Summarize',
          organizationId,
          outputType: 'TEXT',
          toolBase: 'summarize'
        },
        {
          description:
            'Create embeddings from text. This tool supports all languages.',
          inputType: 'TEXT',
          name: 'Create Embeddings',
          organizationId,
          outputType: 'TEXT', // FIXME make this none
          toolBase: 'create-embeddings'
        },
        {
          description:
            'Convert text to speech. This tool supports all languages.',
          inputType: 'TEXT',
          name: 'Text to Speech',
          organizationId,
          outputType: 'AUDIO',
          toolBase: 'text-to-speech'
        }
      ]

      const result = await toolRepository.createMany(defaultTools)
      result.data.forEach((tool) => {
        emitToolMutationEvent(tool, websocketsService)
      })
      return result
    }
  }
}

const emitToolMutationEvent = (
  entity: ToolEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, TOOL_ENTITY_KEY]
  })
}

export type ToolsService = ReturnType<typeof createToolsService>
