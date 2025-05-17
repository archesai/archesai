import type { Server as HttpServer } from 'http'
import type { Socket } from 'socket.io'

import { Server as WebsocketsServer } from 'socket.io'

import type { ConfigService } from '#config/config.service'

import { Logger } from '#logging/logger'
import { RedisIoAdapter } from '#websockets/adapters/redis-io.adapter'

type ArchesWebsocketsSocket = Socket<never, never, never, { sub: string }>

export class WebsocketsService {
  public io?: WebsocketsServer
  public get ioServer(): undefined | WebsocketsServer {
    return this.io
  }
  private readonly configService: ConfigService
  private readonly logger = new Logger(WebsocketsService.name)

  constructor(configService: ConfigService) {
    this.configService = configService
  }

  public broadcastEvent(room: string, event: string, data: unknown): void {
    this.io?.to(room).emit(event, data)
  }

  public async handleConnection(socket: ArchesWebsocketsSocket) {
    try {
      const token = this.getTokenFromSocket(socket)
      // const { sub } = this.jwtService.verify<AccessTokenDecodedJwt>(token)
      const sub = token // FIXME

      await socket.join(sub)
      socket.data.sub = sub
      this.logger.log(`websocket connection successful`, {
        room: sub,
        socketId: socket.id,
        sub: sub
      })
    } catch (error) {
      this.logger.error(`websocket connection error`, { error })
      socket.disconnect()
    }
  }

  public handleDisconnect(socket: ArchesWebsocketsSocket, reason: string) {
    this.logger.log(`websocket disconnected`, {
      reason,
      rooms: socket.rooms,
      socketId: socket.id,
      sub: socket.data.sub
    })
  }

  public async setupWebsocketAdapter(httpServer: HttpServer): Promise<void> {
    this.logger.debug('setting up websockets adapter')
    this.io = new WebsocketsServer(httpServer, {
      cors: {
        credentials: true,
        methods: ['GET', 'POST'],
        origin: '*'
      },
      transports: ['websocket']
    })

    if (this.configService.get('redis.enabled')) {
      const redisIoAdapter = new RedisIoAdapter(this.configService, this.logger)
      await redisIoAdapter.connectToRedis()
      this.io.adapter(redisIoAdapter.adapterConstructor!)
      this.logger.debug('redis adapter attached to socket.io')
    }

    this.io.engine.on('connection_error', (error) => {
      this.logger.error('socket engine connection error', { error })
    })

    this.io.on('connection', (socket: ArchesWebsocketsSocket) => {
      ;(async () => {
        await this.handleConnection(socket)
        socket.on('disconnect', (reason) => {
          this.handleDisconnect(socket, reason)
        })
      })().catch((error: unknown) => {
        this.logger.error('error in connection handler', { error })
      })
    })

    this.logger.debug('websockets adapter setup complete')
  }

  private getTokenFromSocket(socket: ArchesWebsocketsSocket) {
    const cookie = socket.handshake.headers.cookie
    if (!cookie) {
      throw new Error('No cookie provided in websocket handshake')
    }
    const token = decodeURIComponent(cookie)
      .split(';')
      .find((c) => c.trim().startsWith('archesai.accessToken='))
      ?.split('=')[1]

    if (!token) {
      throw new Error('Invalid cookie provided in websocket handshake')
    }
    return token
  }
}
