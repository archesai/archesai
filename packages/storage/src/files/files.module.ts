import type { ModuleMetadata } from '@archesai/core'

import {
  createModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { FilesController } from '#files/files.controller'
import { FilesService } from '#files/files.service'
import { StorageModule } from '#storage/storage.module'
import { StorageService } from '#storage/storage.service'

export const FilesModuleDefinition: ModuleMetadata = {
  imports: [StorageModule, WebsocketsModule],
  providers: [
    {
      inject: [StorageService, WebsocketsService],
      provide: FilesService,
      useFactory: (
        storageService: StorageService,
        websocketsService: WebsocketsService
      ): FilesService => {
        return new FilesService(storageService, websocketsService)
      }
    },
    {
      inject: [FilesService],
      provide: FilesController,
      useFactory: (filesService: FilesService): FilesController => {
        return new FilesController(filesService)
      }
    }
  ]
}

export const FilesModule = (() =>
  createModule(class FilesModule {}, FilesModuleDefinition))()
