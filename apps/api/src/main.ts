import { NestFactory } from '@nestjs/core'
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger'
import { RedisStore } from 'connect-redis'
import cookieParser from 'cookie-parser'
import session from 'express-session'
import { readFileSync } from 'fs'
import helmet from 'helmet'
import { Logger } from 'nestjs-pino'
import passport from 'passport'
import { createClient } from 'redis'

import { AppModule } from '@/src/app.module'
import { RedisIoAdapter } from '@/src/common/adapters/redis-io.adapter'
import { AggregateFieldResult, Metadata } from '@/src/common/dto/paginated.dto'
import { FieldAggregate, FieldFilter } from '@/src/common/dto/search-query.dto'
import { ConfigService } from '@/src/config/config.service'
import { apiReference } from '@scalar/nestjs-api-reference'

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    bufferLogs: true,
    rawBody: true
  })
  const configService = app.get(ConfigService)

  // Docs Setup
  if (configService.get('server.docs.enabled')) {
    const swaggerConfig = new DocumentBuilder()
      .setTitle('Arches AI API')
      .setDescription('The Arches AI API')
      .setVersion('v1')
      .addBearerAuth()
      .addCookieAuth()
      .addServer(
        `${configService.get('tls.enabled') ? 'https://' : 'http://'}${configService.get('server.host')}`!
      )
      .build()
    const document = SwaggerModule.createDocument(app, swaggerConfig, {
      extraModels: [FieldFilter, FieldAggregate, AggregateFieldResult, Metadata]
    })
    app.use(
      '/docs',
      apiReference({
        spec: {
          content: document
        },
        theme: 'purple'
      })
    )
  }

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

  // Start listening for requests
  await app.listen(3001)
}

bootstrap()
