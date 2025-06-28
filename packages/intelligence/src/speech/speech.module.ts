import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { SpeechService } from '#speech/speech.service'

export const SpeechModuleDefinition: ModuleMetadata = {
  exports: [SpeechService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: SpeechService,
      useFactory: (configService: ConfigService) =>
        new SpeechService(configService)
    }
  ]
}

export const SpeechModule = (() =>
  createModule(class SpeechModule {}, SpeechModuleDefinition))()
