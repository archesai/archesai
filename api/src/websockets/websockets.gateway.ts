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
    origin: ['https://platform.archesai.com', 'http://localhost:3000', 'http://arches-api:3001']
  },
  transports: ['websocket']
})
export class WebsocketsGateway implements OnGatewayConnection, OnGatewayDisconnect, OnGatewayInit {
  @WebSocketServer()
  server: Server

  private readonly logger: Logger = new Logger('WebsocketsGateway')

  constructor(
    private readonly authService: AuthService,
    private readonly usersService: UsersService,
    private readonly websocketsService: WebsocketsService
  ) {}

  afterInit(server: Server) {
    this.websocketsService.socket = server
    server.on('error', (error) => {
      this.logger.error(`WebSocket error: ${error}`)
    })
  }

  async handleConnection(socket: Socket) {
    try {
      const cookie = socket.handshake.headers.cookie
      if (!cookie) {
        throw new Error('No cookie provided')
      }
      const token = decodeURIComponent(cookie)
        .split(';')
        .find((c) => c.trim().startsWith('archesai.accessToken='))
        ?.split('=')[1]
      if (!token) {
        throw new Error('No jwt token found in cookie')
      }

      // Remove the 's:' prefix if it exists. Express adds this to the beginning when its a signed cookie
      const cleanToken = token.startsWith('s:') ? token.slice(2) : token

      const { sub: id } = await this.authService.verifyToken(cleanToken)
      const user = await this.usersService.findOne(null, id)
      this.logger.debug(`Connected with websockets ${user.defaultOrgname}`)
      socket.join(user.defaultOrgname)
    } catch (error) {
      this.logger.error(error)
      socket.disconnect()
    }
  }

  // Implement the handleDisconnect method
  async handleDisconnect(socket: Socket) {
    this.logger.log(`Disconnected ${socket.id} ${Array.from(socket.rooms).toString()}`)
  }
}
