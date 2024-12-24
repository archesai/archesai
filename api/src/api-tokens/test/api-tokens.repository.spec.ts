import { Test, TestingModule } from '@nestjs/testing'

import { PrismaService } from '../../prisma/prisma.service'
import { ApiTokenRepository } from '../api-token.repository'
import { CreateApiTokenDto } from '../dto/create-api-token.dto'
import { ApiTokenModel, RoleTypeEnum } from '../entities/api-token.entity'

describe('ApiTokenRepository', () => {
  let repository: ApiTokenRepository
  let prismaService: PrismaService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ApiTokenRepository,
        {
          provide: PrismaService,
          useValue: {
            apiToken: {
              create: jest.fn()
            }
          }
        }
      ]
    }).compile()

    repository = module.get<ApiTokenRepository>(ApiTokenRepository)
    prismaService = module.get<PrismaService>(PrismaService)
  })

  describe('create', () => {
    it('should create an API token', async () => {
      // Arrange
      const orgname = 'test-org'
      const createApiTokenDto: CreateApiTokenDto = {
        domains: '*',
        name: 'Test Token',
        role: RoleTypeEnum.USER
      }
      const overrides = {
        id: 'test-id',
        key: 'test-key',
        username: 'test-user',
        orgname
      }

      const expectedResult: ApiTokenModel = {
        ...createApiTokenDto,
        ...overrides,
        createdAt: new Date(),
        updatedAt: new Date()
      }

      prismaService.apiToken.create = jest
        .fn()
        .mockResolvedValue(expectedResult)

      // Act
      const result = await repository.create({
        ...createApiTokenDto,
        ...overrides
      })

      // Assert
      expect(prismaService.apiToken.create).toHaveBeenCalledWith({
        data: {
          ...createApiTokenDto,
          ...overrides
        }
      })
      expect(result).toEqual(expectedResult)
    })
  })
})
