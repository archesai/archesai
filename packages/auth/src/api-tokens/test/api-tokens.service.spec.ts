import type { TestingModule } from '@nestjs/testing'

import { Test } from '@nestjs/testing'

import { ConfigService, WebsocketsService } from '@archesai/core'
import { API_TOKEN_ENTITY_KEY } from '@archesai/schemas'

import { ApiTokenRepository } from '#api-tokens/api-token.repository'
import { ApiTokensService } from '#api-tokens/api-tokens.service'
import { createRandomApiToken } from '#api-tokens/factories/api-token.factory'
import { JwtService } from '#jwt/jwt.service'

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
      const createTokenDto = {
        name: 'test-token',
        role: 'ADMIN'
      } as const
      const overrides = {
        orgname,
        username: 'test-user'
      }
      const mockedApiToken = createRandomApiToken({
        name: createTokenDto.name,
        orgname: overrides.orgname,
        role: createTokenDto.role
      })

      ;(configService.get as jest.Mock).mockImplementation((key: string) => {
        if (key === 'jwt.expiration') return '3600'
        if (key === 'jwt.secret') return 'secret'
        return null
      })
      ;(jwtService.sign as jest.Mock).mockReturnValue('token')
      ;(apiTokenRepository.create as jest.Mock).mockResolvedValue(
        mockedApiToken
      )

      const result = await service.create({
        name: createTokenDto.name,
        orgname: overrides.orgname,
        role: createTokenDto.role
      })
      expect(result.orgname).toEqual(orgname)
      expect(configService.get).toHaveBeenCalledWith('jwt.expiration')
      expect(configService.get).toHaveBeenCalledWith('jwt.secret')
      expect(jwtService.sign).toHaveBeenCalledWith(
        {
          id: mockedApiToken.id,
          role: createTokenDto.role,
          ...overrides
        },
        {
          expiresIn: '3600s',
          secret: 'secret'
        }
      )
      expect(apiTokenRepository.create).toHaveBeenCalledWith({
        id: mockedApiToken.id,
        key: '*********token',
        ...overrides
      })
      expect(websocketsService.broadcastEvent).toHaveBeenCalledWith('update', {
        queryKey: ['organizations', result.orgname, API_TOKEN_ENTITY_KEY]
      })
      expect(result).toEqual(
        createRandomApiToken({ ...mockedApiToken, key: 'token' })
      )
    })
  })
})
