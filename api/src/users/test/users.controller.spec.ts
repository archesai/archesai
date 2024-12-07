import { createRandomUser } from '@/prisma/factories/user.factory'
import { CommonModule } from '@/src/common/common.module'
import { createMock, DeepMocked } from '@golevelup/ts-jest'
import { INestApplication } from '@nestjs/common'
import { APP_GUARD } from '@nestjs/core'
import { Test, TestingModule } from '@nestjs/testing'
import request from 'supertest'

import { UpdateUserDto } from '../dto/update-user.dto'
import { UsersController } from '../users.controller'
import { UsersService } from '../users.service'

describe('UsersController', () => {
  let app: INestApplication
  let mockedUsersService: DeepMocked<UsersService>

  beforeAll(async () => {
    const moduleRef: TestingModule = await Test.createTestingModule({
      controllers: [UsersController],
      imports: [CommonModule],
      providers: [
        {
          provide: APP_GUARD,
          useValue: createMock({
            canActivate: jest.fn().mockImplementation((context) => {
              const request = context.switchToHttp().getRequest()
              request.user = createRandomUser({
                defaultOrgname: 'test-org',
                id: 'test-id'
              })
              return true
            })
          })
        },
        {
          provide: UsersService,
          useValue: createMock<UsersService>({
            create: jest.fn(),
            deactivate: jest.fn(),
            findAll: jest.fn(),
            findOne: jest.fn(),
            findOneByEmail: jest.fn(),
            remove: jest.fn(),
            setEmailVerified: jest.fn(),
            update: jest.fn()
          })
        }
      ]
    }).compile()

    app = moduleRef.createNestApplication()
    await app.init()

    mockedUsersService = moduleRef.get(UsersService)
  })

  afterAll(async () => {
    await app.close()
  })

  it('should be defined', () => {
    expect(app.get(UsersController)).toBeDefined()
    expect(mockedUsersService).toBeDefined()
  })

  describe('POST /user/deactivate', () => {
    it('should deactivate a user', async () => {
      const response = await request(app.getHttpServer()).post('/user/deactivate').send().expect(201)

      expect(response.body).toEqual({})
      expect(mockedUsersService.deactivate).toHaveBeenCalledWith('test-id')
    })
  })

  describe('GET /user', () => {
    it('should return the current user', async () => {
      const response = await request(app.getHttpServer()).get('/user').expect(200)

      expect(response.body.id).toEqual('test-id')
    })
  })

  describe('PATCH /user', () => {
    it('should update a user', async () => {
      const updateUserDto: UpdateUserDto = {
        firstName: 'John',
        lastName: 'Doe'
      }
      const mockedUser = createRandomUser({
        defaultOrgname: 'test-org',
        id: 'test-id',
        ...updateUserDto
      })
      mockedUsersService.update.mockResolvedValue(mockedUser)

      const response = await request(app.getHttpServer()).patch('/user').send(updateUserDto).expect(200)

      expect(response.body.firstName).toEqual('John')
      expect(response.body.lastName).toEqual('Doe')
      expect(response.body.id).toEqual('test-id')

      expect(mockedUsersService.update).toHaveBeenCalledWith(mockedUser.defaultOrgname, mockedUser.id, updateUserDto)
    })
  })
})
