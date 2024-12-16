import { createRandomMember } from '@/prisma/factories/member.factory'
import { createRandomUser } from '@/prisma/factories/user.factory'
import { CommonModule } from '@/src/common/common.module'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { INestApplication } from '@nestjs/common'
import { APP_GUARD } from '@nestjs/core'
import { Test, TestingModule } from '@nestjs/testing'
import request from 'supertest'

import { MembersController } from '../members.controller'
import { MembersService } from '../members.service'
import { CreateMemberDto } from '../dto/create-member.dto'
import { RoleTypeEnum } from '../entities/member.entity'

describe('MembersController', () => {
  let app: INestApplication
  let mockedMembersService: DeepMocked<MembersService>

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [MembersController],
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
          provide: MembersService,
          useValue: createMock<MembersService>({
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

    mockedMembersService = moduleRef.get(MembersService)
  })

  afterAll(async () => {
    await app.close()
  })

  it('should be defined', () => {
    expect(app.get(MembersController)).toBeDefined()
    expect(mockedMembersService).toBeDefined()
  })

  it('POST /organizations/:orgname/members should validate role', async () => {
    const orgname = 'testOrg'
    const createMemberDto: any = {
      inviteEmail: 'jonathan@gmail.com',
      role: 'BADROLE' as any
    }
    const response = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .set('Content-Type', 'application/json')
      .send(createMemberDto)

    expect(response.status).toBe(400)
  })

  it('POST /organizations/:orgname/members should call service.create', async () => {
    const orgname = 'testOrg'
    const createMemberDto: CreateMemberDto = {
      inviteEmail: 'jonathan@gmail.com',
      role: RoleTypeEnum.ADMIN
    }
    const mockedApiToken = createRandomMember({
      role: RoleTypeEnum.ADMIN
    })
    mockedMembersService.create.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer())
      .post(`/organizations/${orgname}/members`)
      .send(createMemberDto)

    expect(response.status).toBe(201)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })

    expect(mockedMembersService.create).toHaveBeenCalledWith(
      orgname,
      createMemberDto,
      []
    )
  })

  it('GET /organizations/:orgname/members should call service.findAll', async () => {
    const orgname = 'testOrg'
    const mockedApiToken = createRandomMember()
    const mockedPaginatedMembers = {
      aggregates: [],
      metadata: {
        limit: 10,
        offset: 0,
        totalResults: 1
      },
      results: [mockedApiToken]
    }
    mockedMembersService.findAll.mockResolvedValue(mockedPaginatedMembers)

    const response = await request(app.getHttpServer())
      .get(`/organizations/${orgname}/members`)
      .query({})

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedPaginatedMembers,
      results: [
        {
          ...mockedApiToken,
          createdAt: mockedApiToken.createdAt.toISOString(),
          updatedAt: undefined
        }
      ]
    })
    expect(mockedMembersService.findAll).toHaveBeenCalledWith(orgname, {
      aggregates: [],
      endDate: undefined,
      filters: [],
      limit: 10,
      offset: 0,
      sortBy: 'createdAt',
      sortDirection: 'desc',
      startDate: undefined
    })
  })

  it('GET /organizations/:orgname/members/:id should call service.findOne', async () => {
    const orgname = 'testOrg'
    const mockedApiToken = createRandomMember()
    mockedMembersService.findOne.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer()).get(
      `/organizations/${orgname}/members/1`
    )

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })
    expect(mockedMembersService.findOne).toHaveBeenCalledWith('testOrg', '1')
  })

  it('PATCH /organizations/:orgname/members/:id should call service.update', async () => {
    const orgname = 'testOrg'
    const mockedApiToken = createRandomMember()
    mockedMembersService.update.mockResolvedValue(mockedApiToken)

    const response = await request(app.getHttpServer())
      .patch(`/organizations/${orgname}/members/1`)
      .send({ role: RoleTypeEnum.ADMIN })
      .set('Authorization', 'Bearer token')

    expect(response.status).toBe(200)
    expect(response.body).toEqual({
      ...mockedApiToken,
      createdAt: mockedApiToken.createdAt.toISOString(),
      updatedAt: undefined
    })
    expect(mockedMembersService.update).toHaveBeenCalledWith('testOrg', '1', {
      role: RoleTypeEnum.ADMIN
    })
  })

  it('DELETE /organizations/:orgname/members/:id should call service.remove', async () => {
    const response = await request(app.getHttpServer())
      .delete('/organizations/testOrg/members/1')
      .set('Authorization', 'Bearer token')

    expect(response.status).toBe(200)
    expect(response.body).toEqual({})
    expect(mockedMembersService.remove).toHaveBeenCalledWith('testOrg', '1')
  })
})
