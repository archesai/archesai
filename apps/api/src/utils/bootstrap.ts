import type { FastifySchema } from 'fastify'
// eslint-disable-next-line no-restricted-syntax
import type { FastifySerializerCompiler } from 'fastify/types/schema.js'

import fastify from 'fastify'
import fp from 'fastify-plugin'
import {
  serializerCompiler,
  validatorCompiler
} from 'fastify-type-provider-zod'
import qs from 'qs'

import { errorHandlerPlugin } from '@archesai/core'

import { controllersPlugin } from '#plugins/controllers.plugin'
import { corsPlugin } from '#plugins/cors.plugin'
import { docsPlugin } from '#plugins/docs.plugin'
import { healthPlugin } from '#plugins/health.plugin'
import { createContainer } from '#utils/container'

// =================================================================
// UTILITY FUNCTIONS
// =================================================================

export async function bootstrap(): Promise<void> {
  // =================================================================
  // 1. DEPENDENCY INJECTION - Completely functional, no classes!
  // =================================================================
  const container = createContainer()

  // =================================================================
  // 2. CREATE FASTIFY INSTANCE
  // =================================================================
  const app = fastify({
    loggerInstance: container.loggerService.pinoLogger,
    querystringParser: qs.parse,
    trustProxy: true
  })

  app.setValidatorCompiler(validatorCompiler)
  app.setSerializerCompiler(
    serializerCompiler as FastifySerializerCompiler<FastifySchema>
  )

  // Register the centralized error handler
  await app.register(fp(errorHandlerPlugin), {
    includeStack: container.configService.get('logging.level') === 'debug',
    sanitizeHeaders: true
  })

  // =================================================================
  // 3. MIDDLEWARE SETUP
  // =================================================================

  // CORS Configuration
  await app.register(fp(corsPlugin), {
    configService: container.configService
  })

  // Auth Management
  // const sessionsService = app.get(SessionsService)
  // await sessionsService.setup(httpInstance)

  // =================================================================
  // 4. REGISTER ALL FUNCTIONAL PLUGINS
  // =================================================================

  // Docs Setup
  await app.register(fp(docsPlugin), {
    configService: container.configService,
    logger: container.loggerService.logger
  })

  // Register all controllers
  await app.register(controllersPlugin, {
    container
  })

  // Health Check Setup
  await app.register(fp(healthPlugin), {
    databaseService: container.databaseService,
    emailService: container.emailService,
    redisService: container.redisService
  })

  // const gracefulShutdown = () => {
  //   container.loggerService.logger.log('Starting graceful shutdown...')

  //   try {
  //     // await container.databaseService.close()
  //     // await container.redisService.close()
  //     container.loggerService.logger.log('All services closed successfully')
  //   } catch (error) {
  //     container.loggerService.logger.error('Error during shutdown', { error })
  //   }
  // }

  // process.on('SIGTERM', gracefulShutdown)
  // process.on('SIGINT', gracefulShutdown)

  // Websocket Adapter
  await container.websocketsService.setupWebsocketAdapter(app.server)

  await app.listen({
    host: '0.0.0.0',
    port: container.configService.get('server.port')
  })
}
