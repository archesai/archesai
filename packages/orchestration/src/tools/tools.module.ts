import type { ModuleMetadata } from '@archesai/core'
import type { ToolEntity } from '@archesai/domain'

import {
  DatabaseModule,
  DatabaseService,
  LoggingModule,
  Module,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { ToolRepository } from '#tools/tool.repository'
import { ToolsController } from '#tools/tools.controller'
import { ToolsService } from '#tools/tools.service'

const ToolsModuleDefinition: ModuleMetadata = {
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
      useFactory: (databaseService: DatabaseService<ToolEntity>) =>
        new ToolRepository(databaseService)
    }
  ]
}

@Module(ToolsModuleDefinition)
export class ToolsModule {}
