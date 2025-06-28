import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { LlmService } from '#llm/llm.service'

export const LlmModuleDefinition: ModuleMetadata = {
  exports: [LlmService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: LlmService,
      useFactory: (configService: ConfigService) =>
        new LlmService(configService)
    }
  ]
}

export const LlmModule = (() =>
  createModule(class LlmModule {}, LlmModuleDefinition))()
