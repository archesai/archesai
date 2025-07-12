import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

// import type { ArchesApiRequest, ArchesApiResponse } from '@archesai/core'

import { SessionsController } from '#sessions/sessions.controller'
import { SessionsService } from '#sessions/sessions.service'

describe('SessionsController', () => {
  let sessionsController: SessionsController
  // let sessionsenticationService: SessionsService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SessionsController, SessionsController],
      providers: [SessionsService]
    })
      .useMocker(createMock)
      .compile()

    sessionsController = module.get<SessionsController>(SessionsController)
    // sessionsenticationService = module.get<SessionsService>(SessionsService)
  })

  it('should be defined', () => {
    expect(sessionsController).toBeDefined()
  })

  // describe('logout', () => {
  //   it('should logout user', async () => {
  //     const req: ArchesApiRequest = {} as ArchesApiRequest
  //     const res: ArchesApiResponse = {} as ArchesApiResponse

  //     await sessionsController.logout(req, res)
  //     expect(sessionsenticationService.logout).toHaveBeenCalledWith(res)
  //   })
  // })
})
