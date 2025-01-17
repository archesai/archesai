import { Test, TestingModule } from '@nestjs/testing'
import { ConfigService } from '@/src/config/config.service'
// import { z } from 'zod'
import { ConfigModule } from '@/src/config/config.module'
import { RunStatusEnum } from '@/src/runs/entities/run.entity'

// const testConfigSchema = z.object({
//   billing: z.object({
//     enabled: z.coerce.boolean()
//   }),
//   storage: z.object({
//     type: z.string()
//   }),
//   redis: z.object({
//     port: z.coerce.number()
//   })
// })

describe('ConfigModule', () => {
  let module: TestingModule
  let configService: ConfigService

  beforeEach(async () => {
    process.env['ARCHES_STORAGE_TYPE'] = 'minio'
    process.env['ARCHES_STORAGE_ENDPOINT'] = 'http://localhost:9000'
    process.env['ARCHES_BILLING_ENABLED'] = 'true'
    process.env['ARCHES_REDIS_PORT'] = '6379'
    process.env['ARCHES_CONFIG_VALIDATE'] = 'true'

    module = await Test.createTestingModule({
      providers: [ConfigService],
      imports: [ConfigModule]
    }).compile()

    configService = module.get(ConfigService)
  })

  it('should be defined', () => {
    expect(module).toBeDefined()
  })

  it('should have a config service', () => {
    expect(configService).toBeDefined()
  })

  it('should get the whole config', () => {
    expect(configService.getConfig()).toStrictEqual({
      billing: { enabled: true },
      storage: { type: 'minio', endpoint: 'http://localhost:9000' },
      redis: { port: 6379 }
    })
  })

  it('should display its health', () => {
    expect(configService.getHealth()).toMatchObject({
      status: RunStatusEnum.QUEUED
    })
  })

  it('should have a feature billing', () => {
    // it should be a boolean
    expect(configService.get('billing.enabled')).toStrictEqual(true)
  })

  it('should have a storage type', () => {
    expect(configService.get('storage.type')).toStrictEqual('minio')
  })

  it('should convert to number', () => {
    expect(configService.get('redis.port')).toStrictEqual(6379)
  })
})
