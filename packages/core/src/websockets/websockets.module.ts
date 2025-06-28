import type { ModuleMetadata } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { RedisModule } from '#redis/redis.module'
import { RedisService } from '#redis/redis.service'
import { createModule } from '#utils/nest'
import { WebsocketsService } from '#websockets/websockets.service'

export const WebsocketsModuleDefinition: ModuleMetadata = {
  exports: [WebsocketsService],
  imports: [ConfigModule, RedisModule],
  providers: [
    {
      inject: [ConfigService, RedisService],
      provide: WebsocketsService,
      useFactory: (configService: ConfigService, redisService: RedisService) =>
        new WebsocketsService(configService, redisService)
    }
  ]
}

export const WebsocketsModule = (() =>
  createModule(class WebsocketsModule {}, WebsocketsModuleDefinition))()
