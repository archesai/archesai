import { Test, TestingModule } from '@nestjs/testing'
import { HealthController } from './health.controller'
import { createMock } from '@golevelup/ts-jest'

describe('HealthController', () => {
  let controller: HealthController

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [HealthController]
    })
      .useMocker(createMock)
      .compile()

    controller = module.get<HealthController>(HealthController)
  })

  it('should be defined', () => {
    expect(controller).toBeDefined()
  })
})
