import fs from 'node:fs'
import os from 'node:os'
import ospath from 'node:path'
import type { TestingModule } from '@nestjs/testing'

import { Test } from '@nestjs/testing'
import ffmpeg from 'fluent-ffmpeg'

import { FetcherModule } from '@archesai/core'

import { AudioService } from '#audio/audio.service'
import { KeyframesService } from '#keyframes/keyframes.service'

jest.mock('fluent-ffmpeg', () => {
  let mockFfmpeg = jest.fn().mockReturnValue({
    on: jest.fn().mockImplementation(function (event, callback) {
      if (event === 'end') {
        callback()
      }
      return
    }),
    output: jest.fn().mockReturnThis(),
    run: jest.fn(),
    setDuration: jest.fn().mockReturnThis(),
    setStartTime: jest.fn().mockReturnThis()
  })
  mockFfmpeg = jest.fn().mockImplementation((_path, callback) => {
    callback(null, { format: { duration: 100 } })
  })
  return mockFfmpeg
})

jest.mock('fetch')
jest.mock('fs')
jest.mock('fluent-ffmpeg')

describe('AudioService', () => {
  let service: AudioService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AudioService,
        {
          provide: FetcherModule,
          useValue: {
            get: jest.fn(),
            post: jest.fn()
          }
        },
        {
          provide: KeyframesService,
          useValue: {}
        }
      ]
    }).compile()

    service = module.get<AudioService>(AudioService)
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('should trim audio and upload the trimmed file', async () => {
    const audioUrl = 'http://example.com/audio.mp3'
    const startTime = 10
    const duration = 20
    const inputTmpPath = ospath.join(os.tmpdir(), 'original.mp3')
    const outputTmpPath = ospath.join(os.tmpdir(), 'trimmed.mp3')
    const mockResponse = { data: Buffer.from('audio data') }
    const mockUploadUrl = 'http://example.com/trimmed.mp3'

    ;(fs.writeFileSync as jest.Mock).mockImplementation(() => null)
    ;(fs.readFileSync as jest.Mock).mockReturnValue(
      Buffer.from('trimmed audio data')
    )
    ;(fs.statSync as jest.Mock).mockReturnValue({ size: 12345 })
    ;(fs.unlinkSync as jest.Mock).mockImplementation(() => null)

    const result = await service.trimAudio(audioUrl, startTime, duration)

    expect(result).toBe(mockUploadUrl)
    expect(fs.writeFileSync).toHaveBeenCalledWith(
      inputTmpPath,
      mockResponse.data
    )
    expect(ffmpeg.ffprobe).toHaveBeenCalledWith(
      inputTmpPath,
      expect.any(Function)
    )
    expect(ffmpeg).toHaveBeenCalledWith(inputTmpPath)
    expect(fs.unlinkSync).toHaveBeenCalledWith(inputTmpPath)
    expect(fs.unlinkSync).toHaveBeenCalledWith(outputTmpPath)
  })

  it('should return original URL if startTime + duration exceeds audio length', async () => {
    const audioUrl = 'http://example.com/audio.mp3'
    const startTime = 90
    const duration = 20
    const inputTmpPath = ospath.join(os.tmpdir(), 'original.mp3')
    const mockResponse = { data: Buffer.from('audio data') }

    ;(fs.writeFileSync as jest.Mock).mockImplementation(() => null)
    ;(ffmpeg.ffprobe as jest.Mock).mockImplementation((_path, callback) => {
      callback(null, { format: { duration: 100 } })
    })
    ;(fs.unlinkSync as jest.Mock).mockImplementation(() => null)

    const result = await service.trimAudio(audioUrl, startTime, duration)

    expect(result).toBe(audioUrl)
    expect(fs.writeFileSync).toHaveBeenCalledWith(
      inputTmpPath,
      mockResponse.data
    )
    expect(ffmpeg.ffprobe).toHaveBeenCalledWith(
      inputTmpPath,
      expect.any(Function)
    )
    expect(fs.unlinkSync).toHaveBeenCalledWith(inputTmpPath)
  })
})
