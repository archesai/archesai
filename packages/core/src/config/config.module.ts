import type { ModuleMetadata } from '#utils/nest'

import { ConfigController } from '#config/config.controller'
import { ConfigLoader } from '#config/config.loader'
import { ConfigService } from '#config/config.service'
import { ArchesConfigSchema } from '#config/schemas/config.schema'
import { createModule } from '#utils/nest'

export const ConfigModuleDefinition: ModuleMetadata = {
  exports: [ConfigService],
  providers: [
    {
      inject: [ConfigService],
      provide: ConfigController,
      useFactory: (configService: ConfigService) =>
        new ConfigController(configService)
    },
    {
      provide: ConfigService,
      useFactory: () => {
        const config = ConfigLoader.load(ArchesConfigSchema)
        return new ConfigService(config)
      }
    }
  ]
}

export const ConfigModule = (() =>
  createModule(class ConfigModule {}, ConfigModuleDefinition, true))()
