import type { WebsocketsService } from '@archesai/core'
import type { RunEntity, StatusType } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { RUN_ENTITY_KEY } from '@archesai/schemas'

import type { RunRepository } from '#runs/run.repository'

/**
 * Service for runs.
 */
export class RunsService extends BaseService<RunEntity> {
  // private readonly artifactsService: ArtifactsService
  // private readonly flowProducer: FlowProducer
  private readonly runRepository: RunRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    runRepository: RunRepository,
    websocketsService: WebsocketsService
  ) {
    super(runRepository)
    this.runRepository = runRepository
    this.websocketsService = websocketsService
  }

  // override async create(value: RunInsert) {
  //   if (value.runType === RunTypeEnum.PIPELINE_RUN && !value.pipelineId) {
  //     throw new BadRequestException('Pipeline ID is required for pipeline runs')
  //   } else if (value.runType === RunTypeEnum.TOOL_RUN && !value.toolId) {
  //     throw new BadRequestException('Tool ID is required for tool runs')
  //   }

  //   const { data: content } = await this.artifactsService.findMany({
  //     filter: {
  //       id: {
  //         IN: value.inputs.map((input) => input.id)
  //       }
  //     }
  //   })

  //   const run = await this.runRepository.createPipelineRun({
  //     ...content
  //   })

  //   await this.setInputsOrOutputs(run.id, 'inputs', content)

  //   await this.flowProducer.add({
  //     data: {
  //       toolId: 'extract-text'
  //     },
  //     name: 'extract-text',
  //     queueName: 'tool'
  //   })

  //   this.emitMutationEvent(run)
  //   return run
  // }

  // async setInputsOrOutputs(
  //   runId: string,
  //   type: 'inputs' | 'outputs',
  //   content: ArtifactEntity[]
  // ) {
  //   const run = await this.runRepository.setInputsOrOutputs(
  //     runId,
  //     type,
  //     content
  //   )
  //   const runEntity = await this.findOne(run.id)
  //   this.emitMutationEvent(runEntity)
  //   return runEntity
  // }

  // async setProgress(id: string, progress: number) {
  //   const run = await this.runRepository.update(id, {
  //     progress
  //   })
  //   const runEntity = await this.findOne(run.id)
  //   this.emitMutationEvent(runEntity)
  //   return runEntity
  // }

  public async setRunError(id: string, error: string) {
    const run = await this.runRepository.update(id, {
      error
    })
    const runEntity = await this.findOne(run.id)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  public async setStatus(id: string, status: StatusType) {
    switch (status) {
      case 'COMPLETED':
        await this.runRepository.update(id, {
          completedAt: new Date().toISOString()
        })
        await this.runRepository.update(id, {
          progress: 1
        })
        break
      case 'FAILED':
        await this.runRepository.update(id, {
          completedAt: new Date().toISOString()
        })
        break
      case 'PROCESSING':
        await this.runRepository.update(id, {
          startedAt: new Date().toISOString()
        })
        break
    }
    const run = await this.runRepository.update(id, {
      status
    })
    const runEntity = await this.findOne(run.id)
    this.emitMutationEvent(runEntity)
    return runEntity
  }

  protected emitMutationEvent(entity: RunEntity): void {
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname, RUN_ENTITY_KEY]
    })
  }
}
