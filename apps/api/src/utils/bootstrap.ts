import fastify from 'fastify'
import fp from 'fastify-plugin'
import qs from 'qs'

import { errorHandlerPlugin } from '@archesai/core'

import { controllersPlugin } from '#plugins/controllers.plugin'
import { corsPlugin } from '#plugins/cors.plugin'
import { docsPlugin } from '#plugins/docs.plugin'
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

  // =================================================================
  // 5. HEALTH CHECK & SYSTEM ROUTES
  // =================================================================

  app.get(
    '/health',
    {
      schema: {
        response: {
          200: {
            properties: {
              services: {
                properties: {
                  database: { type: 'string' },
                  email: { type: 'string' },
                  redis: { type: 'string' }
                },
                type: 'object'
              },
              timestamp: { type: 'string' },
              uptime: { type: 'number' }
            },
            type: 'object'
          }
        },
        summary: 'Health check endpoint',
        tags: ['System']
      }
    },
    async () => {
      // Use functional services for health checks
      const dbStatus = await checkServiceHealth('database', () =>
        container.databaseService.ping()
      )
      const redisStatus = await checkServiceHealth('redis', () =>
        container.redisService.ping()
      )
      const emailStatus = await checkServiceHealth('email', () =>
        container.emailService.ping()
      )

      return {
        services: {
          database: dbStatus,
          email: emailStatus,
          redis: redisStatus
        },
        timestamp: new Date().toISOString(),
        uptime: process.uptime()
      }
    }
  )

  // =================================================================
  // 6. ERROR HANDLING
  // =================================================================

  // =================================================================
  // 7. GRACEFUL SHUTDOWN
  // =================================================================

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

async function checkServiceHealth(
  serviceName: string,
  healthCheck: () => Promise<boolean>
): Promise<string> {
  try {
    const isHealthy = await healthCheck()
    return isHealthy ? 'healthy' : 'unhealthy'
  } catch (error) {
    console.error(`Health check failed for ${serviceName}:`, error)
    return 'unhealthy'
  }
}
