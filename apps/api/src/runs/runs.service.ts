import { InjectFlowProducer, InjectQueue } from '@nestjs/bullmq'
import { BadRequestException, Injectable } from '@nestjs/common'
import { RunStatus } from '@prisma/client'
import { FlowProducer, Queue } from 'bullmq'

import { BaseService } from '../common/base.service'
import { ContentService } from '../content/content.service'
import { ContentEntity } from '../content/entities/content.entity'
import { PipelinesService } from '../pipelines/pipelines.service'
import { ToolsService } from '../tools/tools.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { CreateRunDto } from './dto/create-run.dto'
import { RunEntity, RunModel, RunTypeEnum } from './entities/run.entity'
import { RunRepository } from './run.repository'
import { RunJob } from './run.processor'

@Injectable()
export class RunsService extends BaseService<
  RunEntity,
  RunModel,
  RunRepository
> {
  constructor(
    private runRepository: RunRepository,
    private websocketsService: WebsocketsService,
    private pipelinesService: PipelinesService,
    private toolsService: ToolsService,
    @InjectFlowProducer('flow') private readonly flowProducer: FlowProducer,
    @InjectQueue('run') private readonly runQueue: Queue<RunJob>,
    private contentService: ContentService
  ) {
    super(runRepository)
  }

  async create(
    createRunDto: CreateRunDto & {
      orgname: string
    }
  ) {
    if (
      createRunDto.runType === RunTypeEnum.PIPELINE_RUN &&
      !createRunDto.pipelineId
    ) {
      throw new BadRequestException('Pipeline ID is required for pipeline runs')
    } else if (
      createRunDto.runType === RunTypeEnum.TOOL_RUN &&
      !createRunDto.toolId
    ) {
      throw new BadRequestException('Tool ID is required for tool runs')
    }

    const runContent = await this.ensureRunContent(
      createRunDto.orgname,
      createRunDto
    )
    if (createRunDto.runType === RunTypeEnum.PIPELINE_RUN) {
      // Create pipeline run
      const run = await this.runRepository.createPipelineRun(
        createRunDto.orgname,
        {
          contentIds: runContent.map((content) => content.id),
          ...createRunDto
        }
      )
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
      const runEntity = this.toEntity(run)
      this.emitMutationEvent(runEntity)
      return runEntity
    } else {
      // Create tool run
      const run = await this.runRepository.createToolRun(createRunDto.orgname, {
        contentIds: runContent.map((content) => content.id),
        ...createRunDto
      })
      // Set inputs
      await this.setInputsOrOutputs(run.id, 'inputs', runContent)
      // Add to tool queue
      const tool = await this.toolsService.findOne(createRunDto.toolId!)
      await this.runQueue.add(tool.toolBase, runContent, {
        jobId: run.id
      })
      // Return run
      const runEntity = this.toEntity(run)
      this.emitMutationEvent(runEntity)
      return runEntity
    }
  }

  async ensureRunContent(orgname: string, createRunDto: CreateRunDto) {
    const runContent: ContentEntity[] = []
    if (createRunDto.contentIds?.length) {
      for (const contentId of createRunDto.contentIds) {
        runContent.push(await this.contentService.findOne(contentId))
      }
    }
    if (createRunDto.text) {
      runContent.push(
        await this.contentService.create({
          name: 'Input Text',
          text: createRunDto.text,
          labels: [],
          orgname,
          url: null
        })
      )
    }
    if (createRunDto.url) {
      runContent.push(
        await this.contentService.create({
          name: 'Input URL',
          url: createRunDto.url,
          labels: [],
          orgname,
          text: null
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
    const runEntity = this.toEntity(run)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  async setProgress(id: string, progress: number) {
    const run = await this.runRepository.update(id, {
      progress
    })
    const runEntity = this.toEntity(run)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  async setRunError(id: string, error: string) {
    const run = await this.runRepository.update(id, {
      error
    })
    const runEntity = this.toEntity(run)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  async setStatus(id: string, status: RunStatus) {
    switch (status) {
      case 'COMPLETE':
        await this.runRepository.update(id, {
          completedAt: new Date()
        })
        await this.runRepository.update(id, {
          progress: 1
        })
        break
      case 'ERROR':
        await this.runRepository.update(id, {
          completedAt: new Date()
        })
        break
      case 'PROCESSING':
        await this.runRepository.update(id, {
          startedAt: new Date()
        })
        break
    }
    const run = await this.runRepository.update(id, {
      status
    })
    const runEntity = this.toEntity(run)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  protected emitMutationEvent(entity: RunEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'runs']
    })
  }

  protected toEntity(model: RunModel): RunEntity {
    return new RunEntity(model)
  }
}
