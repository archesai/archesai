import type { ModuleMetadata } from '@archesai/core'

import { createModule, FetcherModule, FetcherService } from '@archesai/core'

import { AudioService } from '#audio/audio.service'

export const AudioModuleDefinition: ModuleMetadata = {
  exports: [AudioService],
  imports: [FetcherModule],
  providers: [
    {
      inject: [FetcherService],
      provide: AudioService,
      useFactory: (fetcherService: FetcherService) =>
        new AudioService(fetcherService)
    }
  ]
}

export const AudioModule = (() =>
  createModule(class AudioModule {}, AudioModuleDefinition))()
