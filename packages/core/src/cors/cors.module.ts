import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { CorsService } from '#cors/cors.service'
import { createModule } from '#utils/nest'

export const CorsModuleDefinition: ModuleMetadata = {
  exports: [CorsService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: CorsService,
      useFactory: (configService: ConfigService) =>
        new CorsService(configService)
    }
  ]
}

export const CorsModule = (() =>
  createModule(class CorsModule {}, CorsModuleDefinition))()
