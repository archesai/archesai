import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  FetcherModule,
  FetcherService,
  Module
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

@Module(RunpodModuleDefinition)
export class RunpodModule {}
