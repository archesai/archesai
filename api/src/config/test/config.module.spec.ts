import { Test, TestingModule } from '@nestjs/testing'
import { createMock } from '@golevelup/ts-jest'
import { ConfigModule, ConfigService } from '@nestjs/config'
import Joi from 'joi'

describe('ConfigModule', () => {
  let module: TestingModule
  let configService: ConfigService

  beforeEach(async () => {
    process.env.STORAGE_TYPE = 'minio'
    process.env.MINIO_ENDPOINT = 'http://localhost:9000'
    process.env.FEATURE_BILLING = 'true'
    const configModule = await ConfigModule.forRoot({
      ignoreEnvFile: true,
      isGlobal: true,
      validationSchema: Joi.object({
        STORAGE_TYPE: Joi.string()
          .valid('google-cloud', 'local', 'minio')
          .required(),
        MINIO_ENDPOINT: Joi.string().when('STORAGE_TYPE', {
          is: 'minio',
          otherwise: Joi.optional(),
          then: Joi.required()
        }),
        FEATURE_BILLING: Joi.boolean().required()
      })
    })
    module = await Test.createTestingModule({
      imports: [configModule]
    })
      .useMocker(createMock)
      .compile()

    configService = module.get(ConfigService)
  })

  it('should be defined', () => {
    expect(module).toBeDefined()
  })

  it('should have a config service', () => {
    expect(configService).toBeDefined()
  })

  it('should have a feature billing', () => {
    // it should be a boolean
    expect(configService.get('FEATURE_BILLING')).toStrictEqual(true)
  })

  it('should have a storage type', () => {
    const configService = module.get(ConfigService)
    expect(configService.get('STORAGE_TYPE')).toStrictEqual('minio')
  })
})
