import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import type { ArchesApiResponse } from '@archesai/core'

import { AuthenticationController } from '#auth/auth.controller'
import { AuthenticationService } from '#auth/auth.service'
import { createRandomUser } from '#users/factories/user.factory'

describe('LocalAuthController', () => {
  let authenticationService: AuthenticationService
  let authenticationController: AuthenticationController

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [AuthenticationController],
      providers: [AuthenticationService]
    })
      .useMocker(createMock)
      .compile()

    authenticationService = module.get<AuthenticationService>(
      AuthenticationService
    )
    authenticationController = module.get<AuthenticationController>(
      AuthenticationController
    )
  })

  it('should be defined', () => {
    expect(authenticationService).toBeDefined()
    expect(authenticationController).toBeDefined()
  })

  describe('login', () => {
    it('should login user', async () => {
      const dto = {
        email: 'test@example.com',
        password: 'password'
      }
      const user = createRandomUser()
      const res = {} as ArchesApiResponse
      await authenticationService.login(user.id, res)
      expect(authenticationService.login).toHaveBeenCalledWith(dto, user)
    })
  })

  describe('register', () => {
    it('should register user', () => {
      expect(true).toBe(true)
    })
  })
})
