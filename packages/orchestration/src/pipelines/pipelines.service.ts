import type { WebsocketsService } from '@archesai/core'
import type { PipelineInsertModel } from '@archesai/database'
import type { PipelineEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { PIPELINE_ENTITY_KEY } from '@archesai/schemas'

import type { PipelineRepository } from '#pipelines/pipeline.repository'
import type { ToolsService } from '#tools/tools.service'

/**
 * Service for pipelines.
 */
export class PipelinesService extends BaseService<PipelineEntity> {
  private readonly pipelineRepository: PipelineRepository
  private readonly toolsService: ToolsService
  private readonly websocketsService: WebsocketsService

  constructor(
    pipelineRepository: PipelineRepository,
    toolsService: ToolsService,
    websocketsService: WebsocketsService
  ) {
    super(pipelineRepository)
    this.pipelineRepository = pipelineRepository
    this.websocketsService = websocketsService
    this.toolsService = toolsService
  }

  public override async create(value: PipelineInsertModel) {
    const pipeline = await this.pipelineRepository.create({
      ...value
      // steps: value.steps
    })
    return this.findOne(pipeline.id)
  }

  public async createDefaultPipeline(orgname: string) {
    const tools = await this.toolsService.createDefaultTools(orgname)
    const tool = tools.data.find((t) => t.name == 'Extract Text')
    if (!tool) {
      throw new Error('Could not create default pipeline, no extract text tool')
    }

    // const firstId = randomUUID()
    const pipeline = await this.pipelineRepository.create({
      description:
        'This is a default pipeline for indexing arbitrary documents. It extracts text from the document, creates an image from the text, summarizes the text, creates embeddings from the text, and converts the text to speech.',
      organizationId: orgname
      // steps: [
      // {
      //   id: firstId,
      //   name: 'extract-text',
      //   prerequisites: [],
      //   toolId: 'extract-text'
      // },
      // ...tools.data.map((tool) => ({
      //   id: tool.toolBase,
      //   name: tool.name,
      //   prerequisites: [
      //     {
      //       pipelineStepId: firstId
      //     }
      //   ],
      //   toolId: tool.id
      // }))
      // ]
    })

    return this.findOne(pipeline.id)
  }

  protected emitMutationEvent(entity: PipelineEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, PIPELINE_ENTITY_KEY]
    })
  }
}
