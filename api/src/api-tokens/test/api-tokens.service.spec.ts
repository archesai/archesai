import { createRandomApiToken } from '@/prisma/factories/api-token.factory'
import { ConfigService } from '@nestjs/config'
import { JwtService } from '@nestjs/jwt'
import { Test, TestingModule } from '@nestjs/testing'
import { v4 as uuidv4 } from 'uuid'

import { WebsocketsService } from '../../websockets/websockets.service'
import { ApiTokenRepository } from '../api-token.repository'
import { ApiTokensService } from '../api-tokens.service'
import { CreateApiTokenDto } from '../dto/create-api-token.dto'

jest.mock('uuid', () => ({
  v4: jest.fn()
}))

describe('ApiTokensService', () => {
  let service: ApiTokensService
  let apiTokenRepository: ApiTokenRepository
  let configService: ConfigService
  let jwtService: JwtService
  let websocketsService: WebsocketsService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ApiTokensService,
        {
          provide: ApiTokenRepository,
          useValue: {
            create: jest.fn()
          }
        },
        {
          provide: ConfigService,
          useValue: {
            get: jest.fn()
          }
        },
        {
          provide: JwtService,
          useValue: {
            sign: jest.fn()
          }
        },
        {
          provide: WebsocketsService,
          useValue: {
            socket: {
              emit: jest.fn(),
              to: jest.fn().mockReturnThis()
            }
          }
        }
      ]
    }).compile()

    service = module.get<ApiTokensService>(ApiTokensService)
    apiTokenRepository = module.get<ApiTokenRepository>(ApiTokenRepository)
    configService = module.get<ConfigService>(ConfigService)
    jwtService = module.get<JwtService>(JwtService)
    websocketsService = module.get<WebsocketsService>(WebsocketsService)
  })

  describe('create', () => {
    it('should create an API token', async () => {
      const orgname = 'test-org'
      const createTokenDto: CreateApiTokenDto = {
        domains: '*',
        name: 'test-token',
        role: 'ADMIN'
      }
      const additionalData = {
        username: 'test-user'
      }
      const mockedApiToken = createRandomApiToken(createTokenDto)

      ;(uuidv4 as jest.Mock).mockReturnValue(mockedApiToken.id)
      ;(configService.get as jest.Mock).mockImplementation((key: string) => {
        if (key === 'JWT_API_TOKEN_EXPIRATION_TIME') return '3600'
        if (key === 'JWT_API_TOKEN_SECRET') return 'secret'
      })
      ;(jwtService.sign as jest.Mock).mockReturnValue('token')
      ;(apiTokenRepository.create as jest.Mock).mockResolvedValue(mockedApiToken)

      const result = await service.create(orgname, createTokenDto, additionalData)

      expect(uuidv4).toHaveBeenCalled()
      expect(configService.get).toHaveBeenCalledWith('JWT_API_TOKEN_EXPIRATION_TIME')
      expect(configService.get).toHaveBeenCalledWith('JWT_API_TOKEN_SECRET')
      expect(jwtService.sign).toHaveBeenCalledWith(
        {
          domains: createTokenDto.domains,
          id: mockedApiToken.id,
          orgname,
          role: createTokenDto.role,
          username: additionalData.username
        },
        {
          expiresIn: '3600s',
          secret: 'secret'
        }
      )
      expect(apiTokenRepository.create).toHaveBeenCalledWith(orgname, createTokenDto, {
        id: mockedApiToken.id,
        key: '*********token',
        username: additionalData.username
      })
      expect(websocketsService.socket.to).toHaveBeenCalledWith(orgname)
      expect(websocketsService.socket.emit).toHaveBeenCalledWith('update', {
        queryKey: ['organizations', orgname, 'api-tokens']
      })
      expect(result).toEqual(createRandomApiToken({ ...mockedApiToken, key: 'token' }))
    })
  })
})
