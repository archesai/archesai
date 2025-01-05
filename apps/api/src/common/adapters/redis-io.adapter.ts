import { ArchesConfigService } from '@/src/config/config.service'
import { Logger } from '@nestjs/common'
import { IoAdapter } from '@nestjs/platform-socket.io'
import { createAdapter } from '@socket.io/redis-adapter'
import { readFileSync } from 'fs'
import { createClient } from 'redis'
import { ServerOptions } from 'socket.io'

export class RedisIoAdapter extends IoAdapter {
  private adapterConstructor: ReturnType<typeof createAdapter>

  private readonly logger: Logger = new Logger(RedisIoAdapter.name)
  constructor(
    app: any,
    private readonly configService: ArchesConfigService
  ) {
    super(app)
  }

  async connectToRedis(): Promise<void> {
    const retryStrategyOptions = {
      factor: 2,
      initialDelay: 1000,
      maxRetryDelay: 5000
    }

    const connectAndHandleErrors = async (client: any) => {
      client.on('error', (error: any) => {
        this.logger.error('Redis client error: ' + error)
      })

      let retryCount = 0

      const connectWithRetry = async () => {
        try {
          await client.connect()
        } catch (error) {
          this.logger.error('Redis connection error: ' + error)

          const delay =
            retryStrategyOptions.initialDelay *
            Math.pow(retryStrategyOptions.factor, retryCount)

          if (delay <= retryStrategyOptions.maxRetryDelay) {
            this.logger.warn(
              `Reconnecting to Redis in ${delay}ms (attempt ${retryCount + 1})`
            )
            setTimeout(connectWithRetry, delay)
            retryCount += 1
          } else {
            this.logger.error(
              'Max retry attempts reached. Unable to reconnect.'
            )
          }
        }
      }

      await connectWithRetry()
    }

    const pubClient = createClient({
      password: this.configService.get('redis.auth'),
      url: `redis://${this.configService.get('redis.host')}:${this.configService.get('redis.port')}`,
      ...(this.configService.get('redis.ca')
        ? {
            socket: {
              ca: readFileSync(this.configService.get('redis.ca')!),
              rejectUnauthorized: false,
              tls: true
            }
          }
        : {})
    })
    const subClient = pubClient.duplicate()
    await Promise.all([
      connectAndHandleErrors(pubClient),
      connectAndHandleErrors(subClient)
    ])
    this.adapterConstructor = createAdapter(pubClient, subClient)
  }

  createIOServer(port: number, options?: ServerOptions): any {
    const server = super.createIOServer(port, options)
    server.adapter(this.adapterConstructor)
    return server
  }
}
