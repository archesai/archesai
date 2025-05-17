import type { ModuleMetadata } from '@archesai/core'
import type { RunEntity } from '@archesai/domain'

import { LlmModule, RunpodModule } from '@archesai/ai'
import {
  ConfigService,
  DatabaseModule,
  DatabaseService,
  Module,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'
import { StorageModule } from '@archesai/storage'

import { ContentModule } from '#content/content.module'
import { PipelinesModule } from '#pipelines/pipelines.module'
import { RunProcessor } from '#runs/run.processor'
import { RunRepository } from '#runs/run.repository'
import { RunsController } from '#runs/runs.controller'
import { RunsService } from '#runs/runs.service'
import { ToolsModule } from '#tools/tools.module'

export const RunsModuleDefinition: ModuleMetadata = {
  exports: [RunsService],
  imports: [
    StorageModule,
    ContentModule,
    LlmModule,
    RunpodModule,
    PipelinesModule,
    ToolsModule,
    DatabaseModule,
    WebsocketsModule
  ],
  providers: [
    {
      inject: [RunsService],
      provide: RunsController,
      useFactory: (runsService: RunsService) => new RunsController(runsService)
    },
    {
      inject: [RunRepository, WebsocketsService],
      provide: RunsService,
      useFactory: (
        runRepository: RunRepository,
        websocketsService: WebsocketsService
      ) => new RunsService(runRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: RunRepository,
      useFactory: (databaseService: DatabaseService<RunEntity>) =>
        new RunRepository(databaseService)
    },
    {
      inject: [ConfigService, RunsService],
      provide: RunProcessor,
      useFactory: (configService: ConfigService, runsService: RunsService) => {
        return new RunProcessor(configService, runsService)
      }
    }
  ]
}

@Module(RunsModuleDefinition)
export class RunsModule {}
