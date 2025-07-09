import type { ModuleMetadata } from '@archesai/core'
import type { LabelEntity } from '@archesai/schemas'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { LabelRepository } from '#labels/label.repository'
import { LabelsController } from '#labels/labels.controller'
import { LabelsService } from '#labels/labels.service'

export const LabelsModuleDefinition: ModuleMetadata = {
  exports: [LabelsService],
  imports: [DatabaseModule, WebsocketsModule],
  providers: [
    {
      inject: [LabelRepository, WebsocketsService],
      provide: LabelsService,
      useFactory: (
        labelRepository: LabelRepository,
        websocketsService: WebsocketsService
      ) => new LabelsService(labelRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: LabelRepository,
      useFactory: (databaseService: DatabaseService<LabelEntity>) =>
        new LabelRepository(databaseService)
    },
    {
      inject: [LabelsService],
      provide: LabelsController,
      useFactory: (labelsService: LabelsService) =>
        new LabelsController(labelsService)
    }
  ]
}

export const LabelsModule = (() =>
  createModule(class LabelsModule {}, LabelsModuleDefinition))()
