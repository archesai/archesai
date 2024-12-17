import { Logger } from '@nestjs/common'
import { Test, TestingModule } from '@nestjs/testing'
import { Server, Socket } from 'socket.io'

import { AuthService } from '../../auth/services/auth.service'
import { UsersService } from '../../users/users.service'
import { WebsocketsGateway } from '../websockets.gateway'
import { WebsocketsService } from '../websockets.service'

describe('WebsocketsGateway', () => {
  let gateway: WebsocketsGateway
  let authService: AuthService
  let usersService: UsersService
  let websocketsService: WebsocketsService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        WebsocketsGateway,
        {
          provide: AuthService,
          useValue: {
            verifyToken: jest.fn()
          }
        },
        {
          provide: UsersService,
          useValue: {
            findOne: jest.fn()
          }
        },
        {
          provide: WebsocketsService,
          useValue: {
            socket: null
          }
        }
      ]
    }).compile()

    gateway = module.get<WebsocketsGateway>(WebsocketsGateway)
    authService = module.get<AuthService>(AuthService)
    usersService = module.get<UsersService>(UsersService)
    websocketsService = module.get<WebsocketsService>(WebsocketsService)

    Logger.overrideLogger(false)
  })

  it('should be defined', () => {
    expect(gateway).toBeDefined()
  })

  describe('afterInit', () => {
    it('should set the server and handle errors', () => {
      const server = new Server()
      const error = new Error('Test error')
      jest.spyOn(server, 'on').mockImplementation((event, handler) => {
        if (event === 'error') {
          handler(error)
        }

        return server
      })

      gateway.afterInit(server)

      expect(websocketsService.socket).toBe(server)
    })
  })

  describe('handleConnection', () => {
    it('should handle connection with valid token', async () => {
      const socket = {
        disconnect: jest.fn(),
        handshake: {
          headers: {
            cookie: 'archesai.accessToken=validToken'
          }
        },
        join: jest.fn()
      } as unknown as Socket

      const user = { defaultOrgname: 'testOrg' }
      ;(authService.verifyToken as jest.Mock).mockResolvedValue({
        sub: 'userId'
      })
      ;(usersService.findOne as jest.Mock).mockResolvedValue(user)

      await gateway.handleConnection(socket)

      expect(authService.verifyToken).toHaveBeenCalledWith('validToken')
      expect(usersService.findOne).toHaveBeenCalledWith(null, 'userId')
      expect(socket.join).toHaveBeenCalledWith('testOrg')
    })

    it('should handle connection with no cookie', async () => {
      const socket = {
        disconnect: jest.fn(),
        handshake: {
          headers: {}
        }
      } as unknown as Socket

      await gateway.handleConnection(socket)

      expect(socket.disconnect).toHaveBeenCalled()
    })

    it('should handle connection with no token', async () => {
      const socket = {
        disconnect: jest.fn(),
        handshake: {
          headers: {
            cookie: 'someOtherCookie=value'
          }
        }
      } as unknown as Socket

      await gateway.handleConnection(socket)

      expect(socket.disconnect).toHaveBeenCalled()
    })

    it('should handle connection with invalid token', async () => {
      const socket = {
        disconnect: jest.fn(),
        handshake: {
          headers: {
            cookie: 'archesai.accessToken=invalidToken'
          }
        }
      } as unknown as Socket

      ;(authService.verifyToken as jest.Mock).mockRejectedValue(
        new Error('Invalid token')
      )

      await gateway.handleConnection(socket)

      expect(authService.verifyToken).toHaveBeenCalledWith('invalidToken')
      expect(socket.disconnect).toHaveBeenCalled()
    })
  })

  describe('handleDisconnect', () => {
    it('should log disconnection', async () => {
      const socket = {
        id: 'socketId',
        rooms: new Set(['room1', 'room2'])
      } as unknown as Socket

      await gateway.handleDisconnect(socket)
    })
  })
})
