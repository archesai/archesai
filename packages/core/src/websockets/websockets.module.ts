import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { Module } from '#utils/nest'
import { WebsocketsService } from '#websockets/websockets.service'

export const WebsocketsModuleDefinition: ModuleMetadata = {
  exports: [WebsocketsService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: WebsocketsService,
      useFactory: (configService: ConfigService) =>
        new WebsocketsService(configService)
    }
  ]
}

@Module(WebsocketsModuleDefinition)
export class WebsocketsModule {}
