import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { HealthController } from '#health/health.controller'
import { Module } from '#utils/nest'

export const HealthModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule],
  providers: [
    {
      provide: HealthController,
      useFactory: () => new HealthController()
    }
  ]
}

@Module(HealthModuleDefinition)
export class HealthModule {}
