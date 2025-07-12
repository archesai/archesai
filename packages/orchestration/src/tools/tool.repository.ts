import type { DatabaseService } from '@archesai/core'
import type { ToolInsertModel, ToolSelectModel } from '@archesai/database'
import type { ToolEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { ToolTable } from '@archesai/database'
import { ToolEntitySchema } from '@archesai/schemas'

/**
 * Repository for tools.
 */
export class ToolRepository extends BaseRepository<
  ToolEntity,
  ToolInsertModel,
  ToolSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      ToolEntity,
      ToolInsertModel,
      ToolSelectModel
    >
  ) {
    super(databaseService, ToolTable, ToolEntitySchema)
  }

  public async createDefaultTools(organizationId: string): Promise<{
    count: number
    data: ToolEntity[]
  }> {
    return this.createMany([
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
    ])
  }
}
