import type { ModuleMetadata } from '@archesai/core'
import type { ArtifactEntity } from '@archesai/schemas'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'
import { StorageModule, StorageService } from '@archesai/storage'

import { ArtifactRepository } from '#artifacts/artifact.repository'
import { ArtifactsController } from '#artifacts/artifacts.controller'
import { ArtifactsService } from '#artifacts/artifacts.service'

export const ArtifactsModuleDefinition: ModuleMetadata = {
  exports: [ArtifactsService],
  imports: [DatabaseModule, StorageModule, WebsocketsModule],
  providers: [
    {
      inject: [ArtifactRepository, StorageService, WebsocketsService],
      provide: ArtifactsService,
      useFactory: (
        artifactRepository: ArtifactRepository,
        storageService: StorageService,
        websocketsService: WebsocketsService
      ) =>
        new ArtifactsService(
          artifactRepository,
          storageService,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: ArtifactRepository,
      useFactory: (databaseService: DatabaseService<ArtifactEntity>) =>
        new ArtifactRepository(databaseService)
    },
    {
      inject: [ArtifactsService],
      provide: ArtifactsController,
      useFactory: (artifactsService: ArtifactsService) =>
        new ArtifactsController(artifactsService)
    }
  ]
}

export const ArtifactsModule = (() =>
  createModule(class ArtifactsModule {}, ArtifactsModuleDefinition))()
