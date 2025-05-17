import type { ModuleMetadata } from '#utils/nest'

import { FetcherService } from '#fetcher/fetcher.service'
import { Module } from '#utils/nest'

export const FetcherModuleDefinition: ModuleMetadata = {
  exports: [FetcherService],
  providers: [
    {
      provide: FetcherService,
      useFactory: () => new FetcherService()
    }
  ]
}

@Module(FetcherModuleDefinition)
export class FetcherModule {}
