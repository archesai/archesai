import { createRandomApiToken } from '@/prisma/factories/api-token.factory'
import { CommonModule } from '@/src/common/common.module'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { INestApplication, Logger } from '@nestjs/common'
import { Test, TestingModule } from '@nestjs/testing'
import request from 'supertest'

import { ApiTokensController } from '../api-tokens.controller'
import { ApiTokensService } from '../api-tokens.service'
import { CreateApiTokenDto } from '../dto/create-api-token.dto'
import { RoleTypeEnum } from '../entities/api-token.entity'
import { AuthenticatedGuard } from '@/src/auth/guards/authenticated.guard'
import { createRandomUser } from '@/prisma/factories/user.factory'
import { ConfigModule } from '@/src/config/config.module'

describe('ApiTokensController', () => {
  let app: INestApplication
  let mockedApiTokensService: DeepMocked<ApiTokensService>
  let orgname: string
  let username: string

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [ApiTokensController],
      imports: [CommonModule, ConfigModule],
      providers: [
        {
          provide: ApiTokensService,
          useValue: createMock<ApiTokensService>()
        }
      ]
    })
      .overrideGuard(AuthenticatedGuard)
      .useValue({
        canActivate(ctx) {
          const request = ctx.switchToHttp().getRequest()
          request.user = mockUserEntity
          return true
        }
      } as AuthenticatedGuard)
      .compile()
    app = moduleRef.createNestApplication()
    app.useLogger(app.get(Logger))

    await app.init()

    mockedApiTokensService = moduleRef.get(ApiTokensService)

    const mockUserEntity = createRandomUser()
    orgname = mockUserEntity.defaultOrgname
    username = mockUserEntity.username
  })

  afterAll(async () => {
    await app.close()
  })

  it('should be defined', () => {
    expect(app.get(ApiTokensController)).toBeDefined()
    expect(mockedApiTokensService).toBeDefined()
  })

  it('POST /organizations/:orgname/api-tokens should call service.create', async () => {
    const createApiTokenDto: CreateApiTokenDto = {
      domains: '*',
      name: 'testToken',
      role: RoleTypeEnum.ADMIN
    }
    const mockedApiToken = createRandomApiToken(createApiTokenDto)
    mockedApiTokensService.create.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/api-tokens`)
      .send(createApiTokenDto)

    expect(response.status).toBe(201)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })

    expect(mockedApiTokensService.create).toHaveBeenCalledWith({
      domains: '*',
      name: 'testToken',
      role: RoleTypeEnum.ADMIN,
      orgname,
      username
    })
  })

  it('GET /organizations/:orgname/api-tokens should call service.findAll', async () => {
    const mockedApiToken = createRandomApiToken()
    const mockedPaginatedApiTokens = {
      aggregates: [],
      metadata: {
        limit: 10,
        offset: 0,
        totalResults: 1
      },
      results: [mockedApiToken]
    }
    mockedApiTokensService.findAll.mockResolvedValue(mockedPaginatedApiTokens)

    const response = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/api-tokens`)
      .query({})

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedPaginatedApiTokens,
      results: [
        {
          ...mockedApiToken,
          createdAt: mockedApiToken.createdAt.toISOString(),
          updatedAt: undefined
        }
      ]
    })
    expect(mockedApiTokensService.findAll).toHaveBeenCalledWith({
      aggregates: [],
      endDate: undefined,
      filters: [
        {
          field: 'orgname',
          operator: 'equals',
          value: orgname
        }
      ],
      limit: 10,
      offset: 0,
      sortBy: 'createdAt',
      sortDirection: 'desc',
      startDate: undefined
    })
  })

  it('GET /organizations/:orgname/api-tokens/:id should call service.findOne', async () => {
    const mockedApiToken = createRandomApiToken()
    mockedApiTokensService.findOne.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer()).get(
      `/organizations/${orgname}/api-tokens/1`
    )

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })
    expect(mockedApiTokensService.findOne).toHaveBeenCalledWith('1')
  })

  it('PATCH /organizations/:orgname/api-tokens/:id should call service.update', async () => {
    const mockedApiToken = createRandomApiToken()
    mockedApiTokensService.update.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer())
      .patch(`/organizations/${orgname}/api-tokens/1`)
      .send({ name: 'updatedToken' })
      .set('Authorization', 'Bearer token')

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })
    expect(mockedApiTokensService.update).toHaveBeenCalledWith('1', {
      name: 'updatedToken'
    })
  })

  it('DELETE /organizations/:orgname/api-tokens/:id should call service.remove', async () => {
    const response = await request(app.getHttpServer())
      .delete(`/organizations/${orgname}/api-tokens/1`)
      .set('Authorization', 'Bearer token')

    expect(response.status).toBe(200)
    expect(response.body).toEqual({})
    expect(mockedApiTokensService.remove).toHaveBeenCalledWith('1')
  })
})
