import { readFileSync } from 'node:fs'
import type { RedisClientType } from 'redis'

import { createAdapter } from '@socket.io/redis-adapter'
import { createClient } from 'redis'

import type { ConfigService } from '#config/config.service'
import type { Logger } from '#logging/logger'

export class RedisIoAdapter {
  public adapterConstructor?: ReturnType<typeof createAdapter>
  private readonly configService: ConfigService
  private readonly logger: Logger

  constructor(configService: ConfigService, logger: Logger) {
    this.configService = configService
    this.logger = logger
  }

  public async connectToRedis(): Promise<void> {
    if (!this.configService.get('redis.enabled')) {
      return
    }

    const retryStrategyOptions = {
      factor: 2,
      initialDelay: 1000,
      maxRetryDelay: 5000
    }

    const connectAndHandleErrors = async (client: RedisClientType) => {
      client.on('error', (error: unknown) => {
        this.logger.error('redis client error', { error })
      })

      let retryCount = 0

      const connectWithRetry = async () => {
        try {
          await client.connect()
        } catch (error) {
          this.logger.error('redis connection error', { error })

          const delay =
            retryStrategyOptions.initialDelay *
            Math.pow(retryStrategyOptions.factor, retryCount)

          if (delay <= retryStrategyOptions.maxRetryDelay) {
            this.logger.warn(
              `reconnecting to redis in ${delay.toString()}ms (attempt ${(retryCount + 1).toString()})`
            )
            setTimeout(() => {
              connectWithRetry().catch((error: unknown) => {
                this.logger.error('redis connection error', { error })
              })
            }, delay)
            retryCount += 1
          } else {
            this.logger.error(
              'redis connection error - max retry attempts reached'
            )
          }
        }
      }

      await connectWithRetry()
    }

    const redisCa = this.configService.get('redis.ca')
    const redisHost = this.configService.get('redis.host')
    const redisPort = this.configService.get('redis.port').toString()
    const redisAuth = this.configService.get('redis.auth')
    const pubClient = createClient({
      ...(redisAuth ? { password: redisAuth } : {}),
      url: `redis://${redisHost}:${redisPort}`,
      ...(redisCa ?
        {
          socket: {
            ca: readFileSync(redisCa),
            host: redisHost,
            rejectUnauthorized: false,
            tls: true
          }
        }
      : {})
    })
    const subClient = pubClient.duplicate()
    await Promise.all([
      connectAndHandleErrors(pubClient as RedisClientType),
      connectAndHandleErrors(subClient as RedisClientType)
    ])
    this.adapterConstructor = createAdapter(pubClient, subClient)
  }
}
