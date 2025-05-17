import type { ModuleMetadata } from '@archesai/core'
import type { ContentEntity } from '@archesai/domain'

import {
  DatabaseModule,
  DatabaseService,
  Module,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'
import { StorageModule, StorageService } from '@archesai/storage'

import { ContentController } from '#content/content.controller'
import { ContentRepository } from '#content/content.repository'
import { ContentService } from '#content/content.service'

export const ContentModuleDefinition: ModuleMetadata = {
  exports: [ContentService],
  imports: [DatabaseModule, StorageModule, WebsocketsModule],
  providers: [
    {
      inject: [ContentRepository, StorageService, WebsocketsService],
      provide: ContentService,
      useFactory: (
        contentRepository: ContentRepository,
        storageService: StorageService,
        websocketsService: WebsocketsService
      ) =>
        new ContentService(contentRepository, storageService, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: ContentRepository,
      useFactory: (databaseService: DatabaseService<ContentEntity>) =>
        new ContentRepository(databaseService)
    },
    {
      inject: [ContentService],
      provide: ContentController,
      useFactory: (contentService: ContentService) =>
        new ContentController(contentService)
    }
  ]
}

@Module(ContentModuleDefinition)
export class ContentModule {}
