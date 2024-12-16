import { InjectFlowProducer, InjectQueue } from '@nestjs/bullmq'
import { BadRequestException, Injectable, Logger } from '@nestjs/common'
import { RunStatus } from '@prisma/client'
import { FlowProducer, Queue } from 'bullmq'

import { BaseService } from '../common/base.service'
import { ContentService } from '../content/content.service'
import { ContentEntity } from '../content/entities/content.entity'
import { PipelinesService } from '../pipelines/pipelines.service'
import { ToolsService } from '../tools/tools.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { CreateRunDto } from './dto/create-run.dto'
import { RunEntity, RunModel } from './entities/run.entity'
import { RunRepository } from './run.repository'

@Injectable()
export class RunsService extends BaseService<
  RunEntity,
  CreateRunDto,
  any,
  RunRepository,
  RunModel
> {
  private logger = new Logger(RunsService.name)

  constructor(
    private runRepository: RunRepository,
    private websocketsService: WebsocketsService,
    private pipelinesService: PipelinesService,
    private toolsService: ToolsService,
    @InjectFlowProducer('flow') private readonly flowProducer: FlowProducer,
    @InjectQueue('run') private readonly runQueue: Queue,
    private contentService: ContentService
  ) {
    super(runRepository)
  }

  async create(orgname: string, createRunDto: CreateRunDto) {
    if (createRunDto.runType === 'PIPELINE_RUN' && !createRunDto.pipelineId) {
      throw new BadRequestException('Pipeline ID is required for pipeline runs')
    } else if (createRunDto.runType === 'TOOL_RUN' && !createRunDto.toolId) {
      throw new BadRequestException('Tool ID is required for tool runs')
    }

    const runContent = await this.ensureRunContent(orgname, createRunDto)
    if (createRunDto.runType === 'PIPELINE_RUN') {
      // Create pipeline run
      const run = await this.runRepository.createPipelineRun(orgname, {
        contentIds: runContent.map((content) => content.id),
        ...createRunDto
      })
      // Set inputs
      await this.setInputsOrOutputs(run.id, 'inputs', runContent)
      // Add to flow queue
      await this.flowProducer.add({
        data: {
          toolId: 'extract-text'
        },
        name: 'extract-text',
        queueName: 'tool'
      })
      // Return run
      this.emitMutationEvent(orgname)
      return this.toEntity(run)
    } else if (createRunDto.runType === 'TOOL_RUN') {
      // Create tool run
      const run = await this.runRepository.createToolRun(orgname, {
        contentIds: runContent.map((content) => content.id),
        ...createRunDto
      })
      // Set inputs
      await this.setInputsOrOutputs(run.id, 'inputs', runContent)
      // Add to tool queue
      const tool = await this.toolsService.findOne(orgname, createRunDto.toolId)
      await this.runQueue.add(
        tool.toolBase,
        {
          inputs: runContent
        },
        {
          jobId: run.id
        }
      )
      // Return run
      this.emitMutationEvent(orgname)
      return this.toEntity(run)
    }
  }

  async ensureRunContent(orgname: string, createRunDto: CreateRunDto) {
    const runContent: ContentEntity[] = []
    if (createRunDto.contentIds?.length) {
      for (const contentId of createRunDto.contentIds) {
        runContent.push(await this.contentService.findOne(orgname, contentId))
      }
    }
    if (createRunDto.text) {
      runContent.push(
        await this.contentService.create(orgname, {
          name: 'Input Text',
          text: createRunDto.text,
          labels: []
        })
      )
    }
    if (createRunDto.url) {
      runContent.push(
        await this.contentService.create(orgname, {
          name: 'Input URL',
          url: createRunDto.url,
          labels: []
        })
      )
    }
    if (!runContent.length) {
      throw new BadRequestException('No input content provided')
    }

    return runContent
  }

  async setInputsOrOutputs(
    runId: string,
    type: 'inputs' | 'outputs',
    content: ContentEntity[]
  ) {
    const run = await this.runRepository.setInputsOrOutputs(
      runId,
      type,
      content
    )
    this.emitMutationEvent(run.orgname)
    return this.toEntity(run)
  }

  async setProgress(id: string, progress: number) {
    const run = await this.runRepository.updateRaw(null, id, {
      progress
    })
    this.emitMutationEvent(run.orgname)
    return this.toEntity(run)
  }

  async setRunError(id: string, error: string) {
    const run = await this.runRepository.updateRaw(null, id, {
      error
    })
    this.emitMutationEvent(run.orgname)
    return this.toEntity(run)
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case 'COMPLETE':
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date()
        })
        await this.runRepository.updateRaw(null, id, {
          progress: 1
        })
        break
      case 'ERROR':
        await this.runRepository.updateRaw(null, id, {
          completedAt: new Date()
        })
        break
      case 'PROCESSING':
        await this.runRepository.updateRaw(null, id, {
          startedAt: new Date()
        })
        break
    }
    const run = await this.runRepository.updateRaw(null, id, {
      status
    })
    this.emitMutationEvent(run.orgname)
    return this.toEntity(run)
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit('update', {
      queryKey: ['organizations', orgname, 'runs']
    })
  }

  protected toEntity(model: RunModel): RunEntity {
    return new RunEntity(model)
  }
}
