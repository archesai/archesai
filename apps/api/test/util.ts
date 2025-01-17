import { RedisIoAdapter } from '@/src/common/adapters/redis-io.adapter'
import { INestApplication } from '@nestjs/common'
import { Test, TestingModule } from '@nestjs/testing'
import 'tsconfig-paths/register'
import { RedisStore } from 'connect-redis'
import cookieParser from 'cookie-parser'
import session from 'express-session'
import { readFileSync } from 'fs'
import helmet from 'helmet'
import { Logger } from 'nestjs-pino'
import passport from 'passport'
import { createClient } from 'redis'
import request from 'supertest'

import { RegisterDto } from '../src/auth/dto/register.dto'
import { CookiesDto } from '../src/auth/dto/token.dto'
import { OrganizationEntity } from '../src/organizations/entities/organization.entity'
import { UserEntity } from '../src/users/entities/user.entity'
import { UsersService } from '../src/users/users.service'
import { AppModule } from '../src/app.module' // This enables path aliasing based on tsconfig.json
import { ConfigService } from '@/src/config/config.service'

export const createApp = async () => {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule]
  }).compile()
  const app = moduleFixture.createNestApplication()

  const configService = app.get(ConfigService)

  //  Setup Logger
  app.useLogger(app.get(Logger))

  // CORS Configuration
  const allowedOrigins = configService.get('server.cors.origins').split(',')
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
  const sessionSecret = configService.get('session.secret')
  if (!sessionSecret) {
    throw new Error('SESSION_SECRET is not defined')
  }
  const redisClient = createClient({
    password: configService.get('redis.auth'),
    url: `redis://${configService.get('redis.host')}:${configService.get('redis.port')}`,
    ...(configService.get('redis.ca')
      ? {
          socket: {
            ca: readFileSync(configService.get('redis.ca')!),
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
        secure: configService.get('tls.enabled')
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
): Promise<CookiesDto> => {
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
