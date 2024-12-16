import { RedisIoAdapter } from '@/src/common/adapters/redis-io.adapter'
import { INestApplication } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { Test, TestingModule } from '@nestjs/testing'
import 'tsconfig-paths/register'
import { RedisStore } from 'connect-redis'
import cookieParser from 'cookie-parser'
import session from 'express-session'
import { readFileSync } from 'fs-extra'
import helmet from 'helmet'
import { Logger } from 'nestjs-pino'
import passport from 'passport'
import { createClient } from 'redis'
import request from 'supertest'

import { RegisterDto } from '../src/auth/dto/register.dto'
import { TokenDto } from '../src/auth/dto/token.dto'
import { OrganizationEntity } from '../src/organizations/entities/organization.entity'
import { UserEntity } from '../src/users/entities/user.entity'
import { UsersService } from '../src/users/users.service'
import { AppModule } from './../src/app.module' // This enables path aliasing based on tsconfig.json
import { createMock } from '@golevelup/ts-jest'

export const createApp = async () => {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule],
    providers: [
      {
        provide: ConfigService,
        useValue: createMock<ConfigService>()
      }
    ]
  }).compile()
  const app = moduleFixture.createNestApplication()

  const configService = app.get(ConfigService)
  jest.spyOn(configService, 'get').mockImplementation((key: string) => {
    if (key === 'FEATURE_EMAIL') {
      return 'true'
    } else {
      return process.env[key]
    }
  })

  //  Setup Logger
  app.useLogger(app.get(Logger))

  // CORS Configuration
  const allowedOrigins = configService.get<string>('ALLOWED_ORIGINS').split(',')
  app.enableCors({
    allowedHeaders: ['Authorization', 'Content-Type', 'Accept'],
    credentials: true,
    origin: (origin, callback) => {
      if (
        !origin ||
        allowedOrigins.includes(origin) ||
        allowedOrigins[0] === '*'
      ) {
        callback(null, true)
      } else {
        callback(new Error('Not allowed by CORS'))
      }
    }
  })

  // Security Middlewares
  app.use(helmet())

  // Session Management
  const sessionSecret = configService.get<string>('SESSION_SECRET')
  if (!sessionSecret) {
    throw new Error('SESSION_SECRET is not defined')
  }
  const redisClient = createClient({
    password: configService.get('REDIS_AUTH'),
    url: `redis://${configService.get('REDIS_HOST')}:${configService.get('REDIS_PORT')}`,
    ...(configService.get('REDIS_CA_CERT_PATH')
      ? {
          socket: {
            ca: readFileSync(configService.get('REDIS_CA_CERT_PATH')),
            rejectUnauthorized: false,
            tls: true
          }
        }
      : {})
  })
  redisClient.on('error', (error) => {
    app.get(Logger).error('Redis client error: ' + error)
  })
  redisClient.connect().catch(console.error)
  const redisStore = new RedisStore({
    client: redisClient
  })

  app.use(cookieParser(sessionSecret))
  app.use(
    session({
      cookie: {
        httpOnly: true,
        maxAge: 24 * 60 * 60 * 1000,
        sameSite: 'lax',
        secure: configService.get<string>('NODE_ENV') === 'production'
      },
      resave: false,
      saveUninitialized: false,
      secret: sessionSecret,
      store: redisStore
    })
  )

  // Initialize Passport
  app.use(passport.initialize())
  app.use(passport.session())

  // Websocket Adapter
  const redisIoAdapter = new RedisIoAdapter(app, configService)
  await redisIoAdapter.connectToRedis()
  app.useWebSocketAdapter(redisIoAdapter)

  // Enable Shutdown Hooks
  app.enableShutdownHooks()

  return app
}

export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

// Helper function to register a user and return the API token
export const registerUser = async (
  app: INestApplication,
  registerDto: RegisterDto
): Promise<TokenDto> => {
  const res = await request(app.getHttpServer())
    .post('/auth/register')
    .send(registerDto)
  expect(res.status).toBe(201)
  expect(res.type).toBe('application/json')
  expect(res).toSatisfyApiSpec()
  return res.body
}

export const setEmailVerified = async (app: INestApplication, id: string) => {
  const userService = app.get<UsersService>(UsersService)
  await userService.setEmailVerified(id)
}

// Helper function to get user data
export const getUser = async (
  app: INestApplication,
  accessToken: string
): Promise<UserEntity> => {
  const res = await request(app.getHttpServer())
    .get('/user')
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(200)
  expect(res.body.defaultOrgname).toBeTruthy()
  expect(res).toSatisfyApiSpec()
  return res.body
}

// Helper function to check organization data
export const getOrganization = async (
  app: INestApplication,
  orgname: string,
  accessToken: string
): Promise<OrganizationEntity> => {
  const res = await request(app.getHttpServer())
    .get(`/organizations/${orgname}`)
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(200)
  expect(res).toSatisfyApiSpec()
  return res.body
}

// Helper function to deactivate a user
export const deactivateUser = async (
  app: INestApplication,
  accessToken: string
) => {
  const res = await request(app.getHttpServer())
    .post('/user/deactivate')
    .set('Authorization', `Bearer ${accessToken}`)
  expect(res.status).toBe(201)
}

export function testBaseControllerEndpoints(
  getApp: () => INestApplication,
  baseRoute: string,
  accessToken: string,
  testData: {
    createCases: Array<{
      dto: any
      expectedResponse: any
      expectedStatus: number
      name: string
    }>
    findAllCases: Array<{
      expectedResponse: any
      expectedStatus: number
      name: string
    }>
    findOneCases: Array<{
      expectedResponse: any
      expectedStatus: number
      id: string
      name: string
    }>
    removeCases: Array<{
      expectedResponse: any
      expectedStatus: number
      id: string
      name: string
    }>
    updateCases: Array<{
      dto: any
      expectedResponse: any
      expectedStatus: number
      id: string
      name: string
    }>
  }
) {
  describe('Base Controller Endpoints', () => {
    describe('POST ' + baseRoute, () => {
      testData.createCases.forEach((testCase) => {
        it(testCase.name, async () => {
          const app = getApp() // Get the app instance when the test runs
          await request(app.getHttpServer())
            .post(baseRoute)
            .set(
              'Authorization',
              `${accessToken ? `Bearer ${accessToken}` : ''}`
            )
            .send(testCase.dto)
            .expect(testCase.expectedStatus)
            .expect((res) => {
              expect(res.body).toEqual(testCase.expectedResponse)
            })
            .expect((res) => {
              expect(res).toSatisfyApiSpec()
            })
        })
      })
    })

    describe('GET ' + baseRoute, () => {
      testData.findAllCases.forEach((testCase) => {
        it(testCase.name, async () => {
          const app = getApp() // Get the app instance when the test runs
          await request(app.getHttpServer())
            .get(baseRoute)
            .set(
              'Authorization',
              `${accessToken ? `Bearer ${accessToken}` : ''}`
            )
            .expect(testCase.expectedStatus)
            .expect((res) => {
              expect(res.body).toEqual(testCase.expectedResponse)
            })
            .expect((res) => {
              expect(res).toSatisfyApiSpec()
            })
        })
      })
    })

    describe('GET ' + baseRoute + '/:id', () => {
      testData.findOneCases.forEach((testCase) => {
        it(testCase.name, async () => {
          const app = getApp() // Get the app instance when the test runs
          await request(app.getHttpServer())
            .get(`${baseRoute}/${testCase.id}`)
            .set(
              'Authorization',
              `${accessToken ? `Bearer ${accessToken}` : ''}`
            )
            .expect(testCase.expectedStatus)
            .expect((res) => {
              expect(res.body).toEqual(testCase.expectedResponse)
            })
            .expect((res) => {
              expect(res).toSatisfyApiSpec()
            })
        })
      })
    })

    describe('PUT ' + baseRoute + '/:id', () => {
      testData.updateCases.forEach((testCase) => {
        it(testCase.name, async () => {
          const app = getApp() // Get the app instance when the test runs
          await request(app.getHttpServer())
            .put(`${baseRoute}/${testCase.id}`)
            .set(
              'Authorization',
              `${accessToken ? `Bearer ${accessToken}` : ''}`
            )
            .send(testCase.dto)
            .expect(testCase.expectedStatus)
            .expect((res) => {
              expect(res.body).toEqual(testCase.expectedResponse)
            })
            .expect((res) => {
              expect(res).toSatisfyApiSpec()
            })
        })
      })
    })

    describe('DELETE ' + baseRoute + '/:id', () => {
      testData.removeCases.forEach((testCase) => {
        it(testCase.name, async () => {
          const app = getApp() // Get the app instance when the test runs
          await request(app.getHttpServer())
            .delete(`${baseRoute}/${testCase.id}`)
            .set(
              'Authorization',
              `${accessToken ? `Bearer ${accessToken}` : ''}`
            )
            .expect(testCase.expectedStatus)
            .expect((res) => {
              expect(res.body).toEqual(testCase.expectedResponse)
            })
            .expect((res) => {
              expect(res).toSatisfyApiSpec()
            })
        })
      })
    })
  })
}
