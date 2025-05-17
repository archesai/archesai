import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import type { ArchesApiRequest, ArchesApiResponse } from '@archesai/core'

import { AuthenticationController } from '#auth/auth.controller'
import { AuthenticationService } from '#auth/auth.service'

describe('AuthenticationController', () => {
  let authController: AuthenticationController
  let authenticationService: AuthenticationService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [AuthenticationController, AuthenticationController],
      providers: [AuthenticationService]
    })
      .useMocker(createMock)
      .compile()

    authController = module.get<AuthenticationController>(
      AuthenticationController
    )
    authenticationService = module.get<AuthenticationService>(
      AuthenticationService
    )
  })

  it('should be defined', () => {
    expect(authController).toBeDefined()
  })

  describe('logout', () => {
    it('should logout user', async () => {
      const req: ArchesApiRequest = {} as ArchesApiRequest
      const res: ArchesApiResponse = {} as ArchesApiResponse

      await authController.logout(req, res)
      expect(authenticationService.logout).toHaveBeenCalledWith(res)
    })
  })
})
