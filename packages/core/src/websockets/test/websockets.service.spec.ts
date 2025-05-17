import { Server as HttpServer } from 'http'
import type { TestingModule } from '@nestjs/testing'
import type { Server, Socket } from 'socket.io'

import { Test } from '@nestjs/testing'

import { WebsocketsService } from '#websockets/websockets.service'

describe('WebsocketsService', () => {
  let websocketsService: WebsocketsService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [WebsocketsService]
    }).compile()

    websocketsService = module.get<WebsocketsService>(WebsocketsService)
  })

  it('should be defined', () => {
    expect(websocketsService).toBeDefined()
  })

  it('should have socket property initialized to null', () => {
    expect(websocketsService.io).toBeNull()
  })

  it('should allow setting the socket property', () => {
    const mockServer = {} as Server
    websocketsService.io = mockServer
    expect(websocketsService.io).toBe(mockServer)
  })

  describe('afterInit', () => {
    it('should set the server and handle errors', async () => {
      const server = new HttpServer()
      // const error = new Error('Test error')
      // jest.spyOn(server, 'on').mockImplementation((event, handler) => {
      // if (event === 'error') {
      //   handler(error)
      // }

      //   return server
      // })

      await websocketsService.setupWebsocketAdapter(server)

      expect(websocketsService.io).toBe(server)
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

      // ;(jwtService.verify as jest.Mock).mockResolvedValue({
      //   sub: 'userId'
      // })
      ;(socket.join as jest.Mock).mockReturnValue(socket)
      await websocketsService.handleConnection(socket)

      expect(socket.join).toHaveBeenCalledWith('testOrg')
    })

    it('should handle connection with no cookie', async () => {
      const socket = {
        disconnect: jest.fn(),
        handshake: {
          headers: {}
        }
      } as unknown as Socket

      await websocketsService.handleConnection(socket)

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

      await websocketsService.handleConnection(socket)

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

      // ;(jwtService.verify as jest.Mock).mockRejectedValue(
      //   new Error('Invalid token')
      // )

      await websocketsService.handleConnection(socket)

      expect(socket.disconnect).toHaveBeenCalled()
    })
  })

  describe('handleDisconnect', () => {
    it('should log disconnection', () => {
      const socket = {
        data: {
          username: 'testUser'
        },
        id: 'socketId',
        rooms: new Set(['room1', 'room2'])
      } as Socket

      websocketsService.handleDisconnect(socket, 'testReason')
      expect(websocketsService.io).toBeNull()
    })
  })
})
