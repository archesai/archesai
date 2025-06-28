import type { DynamicModule } from '#utils/nest'

import { ConfigModule } from '#config/config.module'
import { ConfigService } from '#config/config.service'
import { RedisService } from '#redis/redis.service'
import { createModule } from '#utils/nest'

const RedisModuleDefinition = {
  // FIXME i dont know if global modules should be done like this. see database module
}

export const RedisModule = (() =>
  createModule(
    class RedisModule {
      public static forRoot(): DynamicModule {
        return {
          exports: [RedisService],
          global: true,
          imports: [ConfigModule],
          module: RedisModule,
          providers: [
            {
              inject: [ConfigService],
              provide: RedisService,
              useFactory: (configService: ConfigService) =>
                new RedisService(configService)
            }
          ]
        }
      }
    },
    RedisModuleDefinition
  ))()
