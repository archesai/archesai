import { createAdapter } from '@socket.io/redis-adapter'

import type { ConfigService } from '#config/config.service'
import type { RedisService } from '#redis/redis.service'

export class RedisIoAdapter {
  public adapterConstructor?: ReturnType<typeof createAdapter>
  private readonly configService: ConfigService
  private readonly redisService: RedisService

  constructor(configService: ConfigService, redisService: RedisService) {
    this.configService = configService
    this.redisService = redisService
  }

  public async connectToRedis(): Promise<void> {
    if (!this.configService.get('redis.enabled')) {
      return
    }

    const { pubClient, subClient } = await this.redisService.createClientPair()
    this.adapterConstructor = createAdapter(pubClient, subClient)
  }
}
