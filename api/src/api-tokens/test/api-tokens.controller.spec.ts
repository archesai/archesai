import { createRandomApiToken } from '@/prisma/factories/api-token.factory'
import { createRandomUser } from '@/prisma/factories/user.factory'
import { CommonModule } from '@/src/common/common.module'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { INestApplication } from '@nestjs/common'
import { APP_GUARD } from '@nestjs/core'
import { Test, TestingModule } from '@nestjs/testing'
import request from 'supertest'

import { ApiTokensController } from '../api-tokens.controller'
import { ApiTokensService } from '../api-tokens.service'
import { CreateApiTokenDto } from '../dto/create-api-token.dto'

describe('ApiTokensController', () => {
  let app: INestApplication
  let mockedApiTokensService: DeepMocked<ApiTokensService>

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [ApiTokensController],
      imports: [CommonModule],
      providers: [
        {
          provide: APP_GUARD,
          useValue: createMock({
            canActivate: jest.fn().mockImplementation((context) => {
              const request = context.switchToHttp().getRequest()
              request.user = createRandomUser({ username: 'testUser' })
              return true
            })
          })
        },
        {
          provide: ApiTokensService,
          useValue: createMock<ApiTokensService>({
            create: jest.fn(),
            findAll: jest.fn(),
            findOne: jest.fn(),
            remove: jest.fn(),
            update: jest.fn()
          })
        }
      ]
    }).compile()

    app = moduleRef.createNestApplication()
    await app.init()

    mockedApiTokensService = moduleRef.get(ApiTokensService)
  })

  afterAll(async () => {
    await app.close()
  })

  it('should be defined', () => {
    expect(app.get(ApiTokensController)).toBeDefined()
    expect(mockedApiTokensService).toBeDefined()
  })

  it('POST /organizations/:orgname/api-tokens should call service.create', async () => {
    const orgname = 'testOrg'
    const createApiTokenDto: CreateApiTokenDto = {
      domains: '*',
      name: 'testToken',
      role: 'ADMIN'
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
      updatedAt: mockedApiToken.updatedAt.toISOString()
    })

    expect(mockedApiTokensService.create).toHaveBeenCalledWith(
      orgname,
      { domains: '*', name: 'testToken', role: 'ADMIN' },
      { username: 'testUser' }
    )
  })

  it('GET /organizations/:orgname/api-tokens should call service.findAll', async () => {
    const orgname = 'testOrg'
    const mockedApiToken = createRandomApiToken()
    const mockedPaginatedApiTokens = {
      aggregates: [],
      metadata: {
        limit: 10,
        offset: 1,
        totalResults: 1
      },
      results: [mockedApiToken]
    }
    mockedApiTokensService.findAll.mockResolvedValue(mockedPaginatedApiTokens)

    const response = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/api-tokens`)
      .query({ limit: 10, offset: 1 })

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedPaginatedApiTokens,
      results: [
        {
          ...mockedApiToken,
          createdAt: mockedApiToken.createdAt.toISOString(),
          updatedAt: mockedApiToken.updatedAt.toISOString()
        }
      ]
    })
    expect(mockedApiTokensService.findAll).toHaveBeenCalledWith(orgname, {
      aggregates: [],
      endDate: undefined,
      filters: [],
      limit: 10,
      offset: 1,
      sortBy: 'createdAt',
      sortDirection: 'desc',
      startDate: undefined
    })
  })

  it('GET /organizations/:orgname/api-tokens/:id should call service.findOne', async () => {
    const orgname = 'testOrg'
    const mockedApiToken = createRandomApiToken()
    mockedApiTokensService.findOne.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer()).get(`/organizations/${orgname}/api-tokens/1`)

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: mockedApiToken.updatedAt.toISOString()
    })
    expect(mockedApiTokensService.findOne).toHaveBeenCalledWith('testOrg', '1')
  })

  it('PATCH /organizations/:orgname/api-tokens/:id should call service.update', async () => {
    const orgname = 'testOrg'
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
      updatedAt: mockedApiToken.updatedAt.toISOString()
    })
    expect(mockedApiTokensService.update).toHaveBeenCalledWith('testOrg', '1', {
      name: 'updatedToken'
    })
  })

  it('DELETE /organizations/:orgname/api-tokens/:id should call service.remove', async () => {
    const response = await request(app.getHttpServer())
      .delete('/organizations/testOrg/api-tokens/1')
      .set('Authorization', 'Bearer token')

    expect(response.status).toBe(200)
    expect(response.body).toEqual({})
    expect(mockedApiTokensService.remove).toHaveBeenCalledWith('testOrg', '1')
  })
})
