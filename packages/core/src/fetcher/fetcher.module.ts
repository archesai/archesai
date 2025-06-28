import type { ModuleMetadata } from '#utils/nest'

import { FetcherService } from '#fetcher/fetcher.service'
import { createModule } from '#utils/nest'

export const FetcherModuleDefinition: ModuleMetadata = {
  exports: [FetcherService],
  providers: [
    {
      provide: FetcherService,
      useFactory: () => new FetcherService()
    }
  ]
}

export const FetcherModule = (() =>
  createModule(class FetcherModule {}, FetcherModuleDefinition))()
