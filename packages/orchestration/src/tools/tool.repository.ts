import type { DatabaseService } from '@archesai/core'

import { BaseRepository } from '@archesai/core'
import { ToolEntity } from '@archesai/domain'

/**
 * Repository for tools.
 */
export class ToolRepository extends BaseRepository<ToolEntity> {
  constructor(databaseService: DatabaseService<ToolEntity>) {
    super(databaseService, 'tools', ToolEntity)
  }

  public async createDefaultTools(orgname: string): Promise<{
    count: number
    data: ToolEntity[]
  }> {
    return this.createMany([
      {
        description:
          'Extract text from a file. This tool supports all file types.',
        inputType: 'TEXT',
        name: 'Extract Text',
        orgname,
        outputType: 'TEXT',
        toolBase: 'extract-text'
      },
      {
        description: 'Create an image from text.',
        inputType: 'TEXT',
        name: 'Text to Image',
        orgname,
        outputType: 'IMAGE',
        toolBase: 'text-to-image'
      },
      {
        description: 'Summarize text. This tool supports all languages.',
        inputType: 'TEXT',
        name: 'Summarize',
        orgname,
        outputType: 'TEXT',
        toolBase: 'summarize'
      },
      {
        description:
          'Create embeddings from text. This tool supports all languages.',
        inputType: 'TEXT',
        name: 'Create Embeddings',
        orgname,
        outputType: 'TEXT', // FIXME make this none
        toolBase: 'create-embeddings'
      },
      {
        description:
          'Convert text to speech. This tool supports all languages.',
        inputType: 'TEXT',
        name: 'Text to Speech',
        orgname,
        outputType: 'AUDIO',
        toolBase: 'text-to-speech'
      }
    ])
  }
}
