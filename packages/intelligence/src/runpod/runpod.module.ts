import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  createModule,
  FetcherModule,
  FetcherService
} from '@archesai/core'

import { RunpodService } from '#runpod/runpod.service'

export const RunpodModuleDefinition: ModuleMetadata = {
  exports: [RunpodService],
  imports: [ConfigModule, FetcherModule],
  providers: [
    {
      inject: [ConfigService, FetcherService],
      provide: RunpodService,
      useFactory: (
        configService: ConfigService,
        fetcherService: FetcherService
      ) => new RunpodService(configService, fetcherService)
    }
  ]
}

export const RunpodModule = (() =>
  createModule(class RunpodModule {}, RunpodModuleDefinition))()
