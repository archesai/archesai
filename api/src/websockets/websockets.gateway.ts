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
import { UsersService } from '../users/users.service'
import { WebsocketsService } from './websockets.service'

@WebSocketGateway({
  connectTimeout: 10000,
  cors: {
    credentials: true,
    origin: [process.env['ARCHES.FRONTEND.HOST']]
  },
  transports: ['websocket']
})
export class WebsocketsGateway
  implements OnGatewayConnection, OnGatewayDisconnect, OnGatewayInit
{
  @WebSocketServer()
  server: Server

  private readonly logger: Logger = new Logger(WebSocketGateway.name)

  constructor(
    private readonly authService: AuthService,
    private readonly usersService: UsersService,
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
      const { sub: id } = await this.authService.verifyToken(token)
      const user = await this.usersService.findOne(id)

      socket.join(user.defaultOrgname)
      socket.data.username = user.username
      this.logger.log(
        `Connected ${user.username} to room ${user.defaultOrgname}`
      )
    } catch (error) {
      this.logger.error(error)
      socket.disconnect()
    }
  }

  async handleDisconnect(socket: Socket) {
    this.logger.log(
      `Disconnected ${socket.data.username} from ${Array.from(socket.rooms).toString()}`
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
