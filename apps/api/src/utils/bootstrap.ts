import fastify from 'fastify'
import fp from 'fastify-plugin'
import qs from 'qs'

import { controllersPlugin } from '#plugins/controllers.plugin'
import { corsPlugin } from '#plugins/cors.plugin'
import { docsPlugin } from '#plugins/docs.plugin'
import { createContainer } from '#utils/container'

export async function bootstrap() {
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

  // Websocket Adapter
  // const websocketsService = app.get(WebsocketsService)
  // await websocketsService.setupWebsocketAdapter(app.getHttpServer())

  await app.listen({
    host: '0.0.0.0',
    port: container.configService.get('server.port')
  })
}

// // =================================================================
// // 5. HEALTH CHECK & SYSTEM ROUTES
// // =================================================================

// app.get(
//   '/health',
//   {
//     schema: {
//       response: {
//         200: {
//           properties: {
//             architecture: { type: 'string' },
//             services: {
//               properties: {
//                 database: { type: 'string' },
//                 email: { type: 'string' },
//                 redis: { type: 'string' }
//               },
//               type: 'object'
//             },
//             status: { type: 'string' },
//             timestamp: { type: 'string' },
//             uptime: { type: 'number' }
//           },
//           type: 'object'
//         }
//       },
//       summary: 'Health check endpoint',
//       tags: ['System']
//     }
//   },
//   async () => {
//     // Use functional services for health checks
//     const dbStatus = await checkServiceHealth('database', () =>
//       container.databaseService.ping()
//     )
//     const redisStatus = await checkServiceHealth('redis', () =>
//       container.redisService.ping()
//     )
//     const emailStatus = await checkServiceHealth('email', () =>
//       container.emailService.ping()
//     )

//     return {
//       architecture: 'Functional Programming + Fastify Plugins',
//       services: {
//         database: dbStatus,
//         email: emailStatus,
//         redis: redisStatus
//       },
//       status: 'healthy',
//       timestamp: new Date().toISOString(),
//       uptime: process.uptime()
//     }
//   }
// )

// =================================================================
// 6. ERROR HANDLING
// =================================================================

//   app.setErrorHandler(async (error, request, reply) => {
//     container.logger.error('Request error', {
//       error: error.message,
//       method: request.method,
//       stack: error.stack,
//       url: request.url
//     })

//     const errorResponse = {
//       error: error.statusCode ? error.message : 'Internal Server Error',
//       path: request.url,
//       statusCode: error.statusCode || 500,
//       timestamp: new Date().toISOString()
//     }

//     reply.status(errorResponse.statusCode).send(errorResponse)
//   })

//   // =================================================================
//   // 7. GRACEFUL SHUTDOWN
//   // =================================================================

//   const gracefulShutdown = async () => {
//     container.logger.info('Starting graceful shutdown...')

//     try {
//       await container.databaseService.close()
//       await container.redisService.close()
//       container.logger.info('All services closed successfully')
//     } catch (error) {
//       container.logger.error('Error during shutdown', { error })
//     }
//   }

//   process.on('SIGTERM', gracefulShutdown)
//   process.on('SIGINT', gracefulShutdown)

//   return app
// }

// =================================================================
// UTILITY FUNCTIONS
// =================================================================

// export async function startFunctionalApp() {
//   try {
//     const app = await createFunctionalApp()

//     const port = process.env.PORT ? parseInt(process.env.PORT, 10) : 3000
//     const host = process.env.HOST || '0.0.0.0'

//     await app.listen({ host, port })

//     console.log(`🚀 Arches AI Functional API is running!`)
//     console.log(`📍 Server: http://${host}:${port}`)
//     console.log(`📚 API Docs: http://${host}:${port}/docs`)
//     console.log(`💚 Health Check: http://${host}:${port}/health`)
//     console.log(`🔧 Architecture: 100% Functional Programming`)
//     console.log(`❌ Classes: ZERO - Completely eliminated!`)
//     console.log(
//       `✨ Benefits: Better performance, easier testing, pure functions`
//     )

//     return app
//   } catch (error) {
//     console.error('❌ Error starting functional application:', error)
//     process.exit(1)
//   }
// }

// // =================================================================
// // APP STARTUP
// // =================================================================

// async function checkServiceHealth(
//   serviceName: string,
//   healthCheck: () => Promise<boolean>
// ): Promise<string> {
//   try {
//     const isHealthy = await healthCheck()
//     return isHealthy ? 'healthy' : 'unhealthy'
//   } catch (error) {
//     console.error(`Health check failed for ${serviceName}:`, error)
//     return 'unhealthy'
//   }
// }
