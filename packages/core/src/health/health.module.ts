import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { HealthController } from '#health/health.controller'
import { createModule } from '#utils/nest'

export const HealthModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule],
  providers: [
    {
      provide: HealthController,
      useFactory: () => new HealthController()
    }
  ]
}

export const HealthModule = (() =>
  createModule(class HealthModule {}, HealthModuleDefinition))()
