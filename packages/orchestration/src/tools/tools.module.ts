import type { ModuleMetadata } from '@archesai/core'
import type { ToolInsertModel, ToolSelectModel } from '@archesai/database'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  LoggingModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { ToolRepository } from '#tools/tool.repository'
import { ToolsController } from '#tools/tools.controller'
import { ToolsService } from '#tools/tools.service'

export const ToolsModuleDefinition: ModuleMetadata = {
  exports: [ToolsService],
  imports: [DatabaseModule, WebsocketsModule, LoggingModule],
  providers: [
    {
      inject: [ToolsService],
      provide: ToolsController,
      useFactory: (toolsService: ToolsService) =>
        new ToolsController(toolsService)
    },
    {
      inject: [ToolRepository, WebsocketsService],
      provide: ToolsService,
      useFactory: (
        toolRepository: ToolRepository,
        websocketsService: WebsocketsService
      ) => new ToolsService(toolRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: ToolRepository,
      useFactory: (
        databaseService: DatabaseService<ToolInsertModel, ToolSelectModel>
      ) => new ToolRepository(databaseService)
    }
  ]
}

export const ToolsModule = (() =>
  createModule(class ToolsModule {}, ToolsModuleDefinition))()
