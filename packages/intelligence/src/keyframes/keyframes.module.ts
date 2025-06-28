import type { ModuleMetadata } from '@archesai/core'

import { createModule } from '@archesai/core'

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

export const KeyframesModule = (() =>
  createModule(class KeyframesModule {}, KeyframesModuleDefinition))()
