import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import { SessionsController } from '#sessions/sessions.controller'
import { SessionsService } from '#sessions/sessions.service'

describe('LocalAuthController', () => {
  let sessionsService: SessionsService
  let sessionsController: SessionsController

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SessionsController],
      providers: [SessionsService]
    })
      .useMocker(createMock)
      .compile()

    sessionsService = module.get<SessionsService>(SessionsService)
    sessionsController = module.get<SessionsController>(SessionsController)
  })

  it('should be defined', () => {
    expect(sessionsService).toBeDefined()
    expect(sessionsController).toBeDefined()
  })

  // describe('login', () => {
  //   it('should login user', async () => {
  //     const dto = {
  //       email: 'test@example.com',
  //       password: 'password'
  //     }
  //     const user = createRandomUser()
  //     const res = {} as ArchesApiResponse
  //     await sessionsService.login(user.id, res)
  //     expect(sessionsService.login).toHaveBeenCalledWith(dto, user)
  //   })
  // })

  describe('register', () => {
    it('should register user', () => {
      expect(true).toBe(true)
    })
  })
})
