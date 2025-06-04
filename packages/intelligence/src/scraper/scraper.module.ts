import type { ModuleMetadata } from '@archesai/core'

import {
  ConfigModule,
  ConfigService,
  FetcherModule,
  FetcherService,
  Module
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

@Module(ScraperModuleDefinition)
export class ScraperModule {}
