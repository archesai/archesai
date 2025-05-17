import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { KeyframesService } from '#keyframes/keyframes.service'

export const KeyframesModuleDefinition: ModuleMetadata = {
  exports: [KeyframesService],
  providers: [
    {
      provide: KeyframesService,
      useFactory: () => new KeyframesService()
    }
  ]
}

@Module(KeyframesModuleDefinition)
export class KeyframesModule {}
