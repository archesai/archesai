import { Test, TestingModule } from '@nestjs/testing'
import { ArchesConfigService } from '../config.service'
import { ConfigModule } from '@nestjs/config'
import { z } from 'zod'
import { loadConfiguration } from '../configuration'

const testConfigSchema = z.object({
  billing: z.object({
    enabled: z.coerce.boolean()
  }),
  storage: z.object({
    type: z.string()
  }),
  redis: z.object({
    port: z.coerce.number()
  })
})

describe('ArchesConfigModule', () => {
  let module: TestingModule
  let configService: ArchesConfigService

  beforeEach(async () => {
    process.env['ARCHES.STORAGE.TYPE'] = 'minio'
    process.env['ARCHES.STORAGE.ENDPOINT'] = 'http://localhost:9000'
    process.env['ARCHES.BILLING.ENABLED'] = 'true'
    process.env['ARCHES.REDIS.PORT'] = '6379'

    module = await Test.createTestingModule({
      providers: [ArchesConfigService],
      imports: [
        ConfigModule.forRoot({
          ignoreEnvFile: true,
          ignoreEnvVars: true,
          load: [() => loadConfiguration(testConfigSchema)]
        })
      ]
    }).compile()

    configService = module.get(ArchesConfigService)
  })

  it('should be defined', () => {
    expect(module).toBeDefined()
  })

  it('should have a config service', () => {
    expect(configService).toBeDefined()
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
