import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import { HealthController } from '#health/health.controller'

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
