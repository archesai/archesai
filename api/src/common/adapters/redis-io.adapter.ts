import { Logger } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { IoAdapter } from '@nestjs/platform-socket.io'
import { createAdapter } from '@socket.io/redis-adapter'
import { readFileSync } from 'fs-extra'
import { createClient } from 'redis'
import { ServerOptions } from 'socket.io'

export class RedisIoAdapter extends IoAdapter {
  private adapterConstructor: ReturnType<typeof createAdapter>

  private readonly logger: Logger = new Logger('RedisIoAdapter')
  constructor(
    app,
    private readonly configService: ConfigService
  ) {
    super(app)
  }

  async connectToRedis(): Promise<void> {
    const retryStrategyOptions = {
      factor: 2,
      initialDelay: 1000,
      maxRetryDelay: 5000
    }

    const connectAndHandleErrors = async (client) => {
      client.on('error', (error) => {
        this.logger.error('Redis client error: ' + error)
      })

      let retryCount = 0

      const connectWithRetry = async () => {
        try {
          await client.connect()
        } catch (error) {
          this.logger.error('Redis connection error: ' + error)

          const delay = retryStrategyOptions.initialDelay * Math.pow(retryStrategyOptions.factor, retryCount)

          if (delay <= retryStrategyOptions.maxRetryDelay) {
            this.logger.warn(`Reconnecting to Redis in ${delay}ms (attempt ${retryCount + 1})`)
            setTimeout(connectWithRetry, delay)
            retryCount += 1
          } else {
            this.logger.error('Max retry attempts reached. Unable to reconnect.')
          }
        }
      }

      await connectWithRetry()
    }

    const pubClient = createClient({
      password: this.configService.get('REDIS_AUTH'),
      url: `redis://${this.configService.get('REDIS_HOST')}:${this.configService.get('REDIS_PORT')}`,
      ...(this.configService.get('REDIS_CA_CERT_PATH')
        ? {
            socket: {
              ca: readFileSync(this.configService.get('REDIS_CA_CERT_PATH')),
              rejectUnauthorized: false,
              tls: true
            }
          }
        : {})
    })
    const subClient = pubClient.duplicate()
    await Promise.all([connectAndHandleErrors(pubClient), connectAndHandleErrors(subClient)])
    this.adapterConstructor = createAdapter(pubClient, subClient)
  }

  createIOServer(port: number, options?: ServerOptions): any {
    const server = super.createIOServer(port, options)
    server.adapter(this.adapterConstructor)
    return server
  }
}
