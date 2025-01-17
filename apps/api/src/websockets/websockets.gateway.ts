import { Logger } from '@nestjs/common'
import {
  OnGatewayConnection,
  OnGatewayDisconnect,
  OnGatewayInit,
  WebSocketGateway,
  WebSocketServer
} from '@nestjs/websockets'
import { Server, Socket } from 'socket.io'

import { AuthService } from '../auth/services/auth.service'
import { WebsocketsService } from './websockets.service'
import { ConfigService } from '@/src/config/config.service'

const configService = new ConfigService()

@WebSocketGateway({
  connectTimeout: 10000,
  cors: {
    credentials: true,
    origin: [configService.get('frontend.host')]
  },
  transports: ['websocket']
})
export class WebsocketsGateway
  implements OnGatewayConnection, OnGatewayDisconnect, OnGatewayInit
{
  @WebSocketServer()
  server: Server

  private readonly logger = new Logger(WebSocketGateway.name)

  constructor(
    private readonly authService: AuthService,
    private readonly websocketsService: WebsocketsService
  ) {}

  afterInit(server: Server) {
    this.websocketsService.socket = server
    server.on('error', (error) => {
      this.logger.error(error)
    })
  }

  async handleConnection(socket: Socket) {
    try {
      const token = await this.getTokenFromSocket(socket)
      const user = await this.authService.verifyToken(token)

      socket.join(user.defaultOrgname)
      socket.data.username = user.username
      this.logger.log(
        {
          socketId: socket.id,
          username: user.username,
          room: user.defaultOrgname
        },
        `websocket connection successful`
      )
    } catch (error) {
      this.logger.error(error, `websocket connection error`)
      socket.disconnect()
    }
  }

  async handleDisconnect(socket: Socket) {
    this.logger.log(
      {
        socketId: socket.id,
        username: socket.data.username,
        rooms: socket.rooms
      },
      `websocket disconnected`
    )
  }

  async getTokenFromSocket(socket: Socket) {
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

    // Remove the 's:' prefix if it exists. Express adds this to the beginning when its a signed cookie
    return token.startsWith('s:') ? token.slice(2) : token
  }
}
