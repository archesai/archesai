import { HttpService } from '@nestjs/axios'
import { Test, TestingModule } from '@nestjs/testing'
import axios from 'axios'
import ffmpeg from 'fluent-ffmpeg'
import * as fs from 'fs'
import * as os from 'os'
import * as ospath from 'path'

import { STORAGE_SERVICE, StorageService } from '../../storage/storage.service'
import { AudioService } from '../audio.service'
import { KeyframesService } from '../keyframes.service'

jest.mock('axios')
jest.mock('fs')
jest.mock('fluent-ffmpeg')

describe('AudioService', () => {
  let service: AudioService
  let storageService: StorageService

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AudioService,
        {
          provide: STORAGE_SERVICE,
          useValue: {
            upload: jest.fn()
          }
        },
        {
          provide: HttpService,
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
    storageService = module.get<StorageService>(STORAGE_SERVICE)
  })

  afterEach(() => {
    jest.clearAllMocks()
  })

  it('should trim audio and upload the trimmed file', async () => {
    const orgname = 'testOrg'
    const audioUrl = 'http://example.com/audio.mp3'
    const startTime = 10
    const duration = 20
    const inputTmpPath = ospath.join(os.tmpdir(), 'original.mp3')
    const outputTmpPath = ospath.join(os.tmpdir(), 'trimmed.mp3')
    const mockResponse = { data: Buffer.from('audio data') }
    const mockUploadUrl = 'http://example.com/trimmed.mp3'

    ;(axios.get as jest.Mock).mockResolvedValue(mockResponse)
    ;(fs.writeFileSync as jest.Mock).mockImplementation(() => {})
    ;(fs.readFileSync as jest.Mock).mockReturnValue(Buffer.from('trimmed audio data'))
    ;(fs.statSync as jest.Mock).mockReturnValue({ size: 12345 })
    ;(fs.unlinkSync as jest.Mock).mockImplementation(() => {})
    ;(ffmpeg as any).mockReturnValue({
      on: jest.fn().mockImplementation(function (event, callback) {
        if (event === 'end') {
          callback()
        }
        return this
      }),
      output: jest.fn().mockReturnThis(),
      run: jest.fn(),
      setDuration: jest.fn().mockReturnThis(),
      setStartTime: jest.fn().mockReturnThis()
    })
    ;(ffmpeg.ffprobe as jest.Mock).mockImplementation((path, callback) => {
      callback(null, { format: { duration: 100 } })
    })
    ;(storageService.upload as jest.Mock).mockResolvedValue(mockUploadUrl)

    const result = await service.trimAudio(orgname, audioUrl, startTime, duration)

    expect(result).toBe(mockUploadUrl)
    expect(axios.get).toHaveBeenCalledWith(audioUrl, {
      responseType: 'arraybuffer'
    })
    expect(fs.writeFileSync).toHaveBeenCalledWith(inputTmpPath, mockResponse.data)
    expect(ffmpeg.ffprobe).toHaveBeenCalledWith(inputTmpPath, expect.any(Function))
    expect(ffmpeg).toHaveBeenCalledWith(inputTmpPath)
    expect(storageService.upload).toHaveBeenCalledWith(orgname, expect.any(String), expect.any(Object))
    expect(fs.unlinkSync).toHaveBeenCalledWith(inputTmpPath)
    expect(fs.unlinkSync).toHaveBeenCalledWith(outputTmpPath)
  })

  it('should return original URL if startTime + duration exceeds audio length', async () => {
    const orgname = 'testOrg'
    const audioUrl = 'http://example.com/audio.mp3'
    const startTime = 90
    const duration = 20
    const inputTmpPath = ospath.join(os.tmpdir(), 'original.mp3')
    const mockResponse = { data: Buffer.from('audio data') }

    ;(axios.get as jest.Mock).mockResolvedValue(mockResponse)
    ;(fs.writeFileSync as jest.Mock).mockImplementation(() => {})
    ;(ffmpeg.ffprobe as jest.Mock).mockImplementation((path, callback) => {
      callback(null, { format: { duration: 100 } })
    })
    ;(fs.unlinkSync as jest.Mock).mockImplementation(() => {})

    const result = await service.trimAudio(orgname, audioUrl, startTime, duration)

    expect(result).toBe(audioUrl)
    expect(axios.get).toHaveBeenCalledWith(audioUrl, {
      responseType: 'arraybuffer'
    })
    expect(fs.writeFileSync).toHaveBeenCalledWith(inputTmpPath, mockResponse.data)
    expect(ffmpeg.ffprobe).toHaveBeenCalledWith(inputTmpPath, expect.any(Function))
    expect(fs.unlinkSync).toHaveBeenCalledWith(inputTmpPath)
  })
})
