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

import { AppModule } from './app.module'
import { RedisIoAdapter } from './common/adapters/redis-io.adapter'
import { AggregateFieldResult, Metadata } from './common/dto/paginated.dto'
import { FieldAggregate, FieldFilter } from './common/dto/search-query.dto'
import { ArchesConfigService } from './config/config.service'

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    bufferLogs: true,
    rawBody: true
  })
  const configService = app.get(ArchesConfigService)

  // Swagger Setup
  if (configService.get('server.swagger.enabled')) {
    const swaggerConfig = new DocumentBuilder()
      .setTitle('Arches AI API')
      .setDescription('The Arches AI API')
      .setVersion('v1')
      .addBearerAuth()
      .addServer(
        `${configService.get('tls.enabled') ? 'https://' : 'http://'}${configService.get('server.host')}`!
      )
      .build()
    const documentFactory = () =>
      SwaggerModule.createDocument(app, swaggerConfig, {
        extraModels: [
          FieldFilter,
          FieldAggregate,
          AggregateFieldResult,
          Metadata
        ]
        // operationIdFactory: (controllerKey: string, methodKey: string) => methodKey
      })

    SwaggerModule.setup('swagger', app, documentFactory, {
      customCss: '.swagger-ui .topbar { display: none }',
      jsonDocumentUrl: 'swagger/json',
      swaggerOptions: {
        persistAuthorization: true,
        tagsSorter: 'alpha'
      },
      yamlDocumentUrl: 'swagger/yaml'
    })
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
