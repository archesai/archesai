import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  createModule,
  FetcherModule,
  FetcherService
} from '@archesai/core'

import { ScraperService } from '#scraper/scraper.service'

export const ScraperModuleDefinition: ModuleMetadata = {
  exports: [ScraperService],
  imports: [ConfigModule, FetcherModule],
  providers: [
    {
      inject: [ConfigService, FetcherService],
      provide: ScraperService,
      useFactory: (
        configService: ConfigService,
        fetcherService: FetcherService
      ) => new ScraperService(configService, fetcherService)
    }
  ]
}

export const ScraperModule = (() =>
  createModule(class ScraperModule {}, ScraperModuleDefinition))()
