import type { ModuleMetadata } from '@archesai/core'

import { createModule } from '@archesai/core'

import { HashingService } from '#hashing/hashing.service'

export const HashingModuleDefinition: ModuleMetadata = {
  exports: [HashingService],
  providers: [
    {
      provide: HashingService,
      useFactory: () => new HashingService()
    }
  ]
}

export const HashingModule = (() =>
  createModule(class HashingModule {}, HashingModuleDefinition))()
