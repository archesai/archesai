import { Test, TestingModule } from '@nestjs/testing'
import { Server } from 'socket.io'

import { WebsocketsService } from '../websockets.service'

describe('WebsocketsService', () => {
  let service: WebsocketsService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [WebsocketsService]
    }).compile()

    service = module.get<WebsocketsService>(WebsocketsService)
  })

  it('should be defined', () => {
    expect(service).toBeDefined()
  })

  it('should have socket property initialized to null', () => {
    expect(service.socket).toBeNull()
  })

  it('should allow setting the socket property', () => {
    const mockServer = {} as Server
    service.socket = mockServer
    expect(service.socket).toBe(mockServer)
  })
})
