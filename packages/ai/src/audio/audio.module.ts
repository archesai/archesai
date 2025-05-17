import type { ModuleMetadata } from '@archesai/core'

import { FetcherModule, FetcherService, Module } from '@archesai/core'

import { AudioService } from '#audio/audio.service'

const AudioModuleDefinition: ModuleMetadata = {
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

@Module(AudioModuleDefinition)
export class AudioModule {}
