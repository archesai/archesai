import { readFileSync } from 'node:fs'
import type { RedisClientType } from '@redis/client'

import { createClient as createRedisClient } from '@redis/client'

import type { ConfigService } from '#config/config.service'
import type { Logger } from '#logging/logger'

export interface RedisConnectionOptions {
  healthCheckInterval?: number
  initialDelay?: number
  maxRetryDelay?: number
  retryAttempts?: number
  retryFactor?: number
}

export const createRedisService = (
  configService: ConfigService,
  logger: Logger,
  options: RedisConnectionOptions = {}
): RedisService => {
  return new RedisService(configService, logger, options)
}

export class RedisService {
  private clients: RedisClientType[] = []
  private readonly configService: ConfigService
  private healthCheckTimer?: NodeJS.Timeout
  private readonly logger: Logger
  private readonly options: Required<RedisConnectionOptions>

  constructor(
    configService: ConfigService,
    logger: Logger,
    options: RedisConnectionOptions = {}
  ) {
    this.logger = logger
    this.configService = configService
    this.options = {
      healthCheckInterval: options.healthCheckInterval ?? 30000,
      initialDelay: options.initialDelay ?? 1000,
      maxRetryDelay: options.maxRetryDelay ?? 30000,
      retryAttempts: options.retryAttempts ?? 10,
      retryFactor: options.retryFactor ?? 2
    }
  }

  public async createClient(): Promise<RedisClientType> {
    if (this.configService.get('redis.mode') === 'disabled') {
      throw new Error('Redis is not enabled')
    }

    const client = this.buildRedisClient()
    await this.connectClientWithRetry(client)
    this.setupClientEventHandlers(client)
    this.clients.push(client)
    return client
  }

  public async createClientPair(): Promise<{
    pubClient: RedisClientType
    subClient: RedisClientType
  }> {
    const pubClient = await this.createClient()
    const subClient = pubClient.duplicate()
    await this.connectClientWithRetry(subClient)
    this.setupClientEventHandlers(subClient)
    this.clients.push(subClient)
    return { pubClient, subClient }
  }

  public async disconnect(): Promise<void> {
    this.stopHealthCheck()
    await Promise.all(
      this.clients.map(async (client) => {
        try {
          if (client.isOpen) {
            await client.close()
          }
        } catch (error) {
          this.logger.error('error disconnecting redis client', { error })
        }
      })
    )
    this.clients = []
  }

  public async healthCheck(): Promise<boolean> {
    try {
      const promises = this.clients.map(async (client) => {
        if (!client.isOpen) {
          return false
        }
        await client.ping()
        return true
      })
      const results = await Promise.all(promises)
      return results.every((result) => result)
    } catch (error) {
      this.logger.error('redis health check failed', { error })
      return false
    }
  }

  public async ping(): Promise<boolean> {
    return this.healthCheck()
  }

  public startHealthCheck(): void {
    if (this.healthCheckTimer) {
      return
    }

    this.healthCheckTimer = setInterval(() => {
      this.healthCheck()
        .then((isHealthy) => {
          if (!isHealthy) {
            this.logger.warn(
              'redis health check failed - some clients may be unhealthy'
            )
          }
        })
        .catch((error: unknown) => {
          this.logger.error('redis health check failed', { error })
        })
    }, this.options.healthCheckInterval)
  }

  public stopHealthCheck(): void {
    if (this.healthCheckTimer) {
      clearInterval(this.healthCheckTimer)
      this.healthCheckTimer = setInterval(() => {
        return
      }, this.options.healthCheckInterval)
    }
  }

  private buildRedisClient(): RedisClientType {
    const redisCa = this.configService.get('redis.ca')
    const redisHost = this.configService.get('redis.host')
    const redisPort = this.configService.get('redis.port').toString()
    const redisAuth = this.configService.get('redis.auth')

    return createRedisClient({
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
  }

  private async connectClientWithRetry(client: RedisClientType): Promise<void> {
    let retryCount = 0

    const connectWithRetry = async (): Promise<void> => {
      try {
        await client.connect()
        this.logger.debug('redis client connected successfully')
      } catch (error) {
        this.logger.error('redis connection error', {
          attempt: retryCount + 1,
          error
        })

        if (retryCount >= this.options.retryAttempts) {
          this.logger.error(
            'redis connection failed - max retry attempts reached'
          )
          throw error
        }

        const delay = Math.min(
          this.options.initialDelay *
            Math.pow(this.options.retryFactor, retryCount),
          this.options.maxRetryDelay
        )

        this.logger.warn(
          `reconnecting to redis in ${delay.toString()}ms (attempt ${(retryCount + 1).toString()}/${this.options.retryAttempts.toString()})`
        )

        await new Promise((resolve) => setTimeout(resolve, delay))
        retryCount += 1
        await connectWithRetry()
      }
    }

    await connectWithRetry()
  }

  private setupClientEventHandlers(client: RedisClientType): void {
    client.on('error', (error: unknown) => {
      this.logger.error('redis client error', { error })
    })

    client.on('connect', () => {
      this.logger.debug('redis client connected')
    })

    client.on('ready', () => {
      this.logger.debug('redis client ready')
    })

    client.on('end', () => {
      this.logger.debug('redis client connection ended')
    })

    client.on('reconnecting', () => {
      this.logger.debug('redis client reconnecting')
    })
  }
}

export type { RedisClientType }
