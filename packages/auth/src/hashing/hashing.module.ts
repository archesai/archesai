import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

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

@Module(HashingModuleDefinition)
export class HashingModule {}
