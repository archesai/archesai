import { Injectable } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import {
  PipelineEntity,
  PipelineWithPipelineStepsModel
} from './entities/pipeline.entity'
import {
  PipelineRepository,
  PipelineStepRepository
} from './pipeline.repository'
import { ToolsService } from '../tools/tools.service'
import { OperatorEnum } from '../common/dto/search-query.dto'

@Injectable()
export class PipelinesService extends BaseService<
  PipelineEntity,
  PipelineWithPipelineStepsModel,
  PipelineRepository
> {
  constructor(
    private pipelineRepository: PipelineRepository,
    private pipelineStepRepository: PipelineStepRepository,
    private websocketsService: WebsocketsService,
    private toolService: ToolsService
  ) {
    super(pipelineRepository)
  }

  async create(data: PipelineEntity) {
    const { pipelineSteps, ...rest } = data
    const pipeline = await this.pipelineRepository.create(rest)
    if (!pipelineSteps) {
      return this.findOne(pipeline.id)
    }
    for (const pipelineStep of pipelineSteps) {
      await this.pipelineStepRepository.create({
        dependsOn: {
          connect: pipelineStep.dependsOn?.map((step) => ({
            id: step.id
          }))
        },
        id: pipelineStep.id,
        name: pipelineStep.name,
        pipelineId: pipeline.id,
        toolId: pipelineStep.toolId
      })
    }

    return this.findOne(pipeline.id)
  }

  async createDefaultPipeline(orgname: string) {
    const pipeline = await this.pipelineRepository.create({
      description:
        'This is a default pipeline for indexing arbitrary documents. It extracts text from the document, creates an image from the text, summarizes the text, creates embeddings from the text, and converts the text to speech.',
      name: 'Default',
      orgname
    })
    const tools = await this.toolService.findAll({
      filters: [
        {
          field: 'orgname',
          operator: OperatorEnum.EQUALS,
          value: orgname
        }
      ]
    })

    // Create first step, this has no dependents
    const firstStep = await this.pipelineStepRepository.create({
      name: 'extract-text',
      pipelineId: pipeline.id,
      toolId: tools.results.find((t) => t.name == 'Extract Text')!.id
    })
    const dependents = tools.results.filter((t) => t.name != 'Extract Text')

    for (const tool of dependents) {
      await this.pipelineStepRepository.create({
        dependsOn: {
          connect: {
            id: firstStep.id
          }
        },
        name: tool.toolBase,
        pipelineId: pipeline.id,
        toolId: tool.id
      })
    }

    return this.toEntity(await this.findOne(pipeline.id))
  }

  async update(id: string, data: Partial<PipelineEntity>) {
    const previousPipeline = await this.pipelineRepository.findOne(id)
    const pipelineStepsToDelete = previousPipeline.pipelineSteps.map(
      (tool) => tool.id
    )

    await this.pipelineRepository.update(id, {
      name: data.name
    })

    await this.pipelineStepRepository.deleteMany({
      filters: [
        {
          field: 'id',
          operator: OperatorEnum.IN,
          value: pipelineStepsToDelete
        }
      ]
    })

    for (const pipelineStep of data.pipelineSteps || []) {
      await this.pipelineStepRepository.create({
        dependsOn: {
          connect: pipelineStep.dependsOn?.map((step) => ({
            id: step.id
          }))
        },
        id: pipelineStep.id,
        name: pipelineStep.name,
        pipelineId: id,
        toolId: pipelineStep.toolId
      })
    }

    return this.findOne(id)
  }

  protected emitMutationEvent(entity: PipelineEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'pipelines']
    })
  }

  protected toEntity(model: PipelineWithPipelineStepsModel): PipelineEntity {
    return new PipelineEntity(model)
  }
}
