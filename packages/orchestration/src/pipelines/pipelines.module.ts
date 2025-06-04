import type { ModuleMetadata } from '@archesai/core'
import type { PipelineEntity } from '@archesai/domain'

import {
  DatabaseModule,
  DatabaseService,
  EventBus,
  EventBusModule,
  Module,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { PipelineRepository } from '#pipelines/pipeline.repository'
import { PipelinesController } from '#pipelines/pipelines.controller'
import { PipelinesService } from '#pipelines/pipelines.service'
import { PipelinesSubscriber } from '#pipelines/pipelines.subscriber'
import { ToolsModule } from '#tools/tools.module'
import { ToolsService } from '#tools/tools.service'

export const PipleinesModuleDefinition: ModuleMetadata = {
  exports: [PipelinesService],
  imports: [DatabaseModule, EventBusModule, ToolsModule, WebsocketsModule],
  providers: [
    {
      inject: [PipelineRepository, ToolsService, WebsocketsService],
      provide: PipelinesService,
      useFactory: (
        pipelineRepository: PipelineRepository,
        toolsService: ToolsService,
        websocketsService: WebsocketsService
      ) =>
        new PipelinesService(
          pipelineRepository,
          toolsService,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: PipelineRepository,
      useFactory: (databaseService: DatabaseService<PipelineEntity>) =>
        new PipelineRepository(databaseService)
    },
    {
      inject: [EventBus, PipelinesService],
      provide: PipelinesSubscriber,
      useFactory: (eventBus: EventBus, pipelinesService: PipelinesService) =>
        new PipelinesSubscriber(eventBus, pipelinesService)
    },
    {
      inject: [PipelinesService],
      provide: PipelinesController,
      useFactory: (pipelinesService: PipelinesService) =>
        new PipelinesController(pipelinesService)
    }
  ]
}

@Module(PipleinesModuleDefinition)
export class PipelinesModule {}
