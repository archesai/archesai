import type { TestingModule } from '@nestjs/testing'

import { createMock } from '@golevelup/ts-jest'
import { Test } from '@nestjs/testing'

import { AudioModule } from '#audio/audio.module'
import { AudioService } from '#audio/audio.service'
import { KeyframesService } from '#keyframes/keyframes.service'

describe('AudioModule', () => {
  let moduleRef: TestingModule

  beforeEach(async () => {
    moduleRef = await Test.createTestingModule({
      imports: [AudioModule]
    })
      .useMocker(createMock)
      .compile()
  })

  it('should resolve exported providers from the ioc container', () => {
    const audioService = moduleRef.get<AudioService>(AudioService)
    const keyframesService = moduleRef.get<KeyframesService>(KeyframesService)
    expect(audioService).toBeDefined()
    expect(keyframesService).toBeDefined()
  })
})
