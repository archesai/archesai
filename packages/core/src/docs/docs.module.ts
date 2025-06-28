import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { DocsService } from '#docs/docs.service'
import { createModule } from '#utils/nest'

export const DocsModuleDefinition: ModuleMetadata = {
  exports: [DocsService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: DocsService,
      useFactory: (configService: ConfigService) =>
        new DocsService(configService)
    }
  ]
}

export const DocsModule = (() =>
  createModule(class DocsModule {}, DocsModuleDefinition))()
