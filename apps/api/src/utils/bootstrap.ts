import fastify from 'fastify'
import fp from 'fastify-plugin'
import {
  serializerCompiler,
  validatorCompiler
} from 'fastify-type-provider-zod'
import qs from 'qs'

import { errorHandlerPlugin } from '@archesai/core'

import { controllersPlugin } from '#app/plugins/controllers.plugin'
import { corsPlugin } from '#app/plugins/cors.plugin'
import { docsPlugin } from '#app/plugins/docs.plugin'
import { healthPlugin } from '#app/plugins/health.plugin'
import { createContainer } from '#utils/container'

export async function bootstrap(): Promise<void> {
  const container = createContainer()

  const app = fastify({
    loggerInstance: container.loggerService.pinoLogger,
    querystringParser: qs.parse,
    trustProxy: true
  })

  app.setValidatorCompiler(validatorCompiler)
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
  app.setSerializerCompiler(serializerCompiler)

  await app.register(fp(errorHandlerPlugin), {
    includeStack: container.configService.get('logging.level') === 'debug',
    sanitizeHeaders: true
  })

  await app.register(fp(corsPlugin), {
    configService: container.configService
  })

  await app.register(fp(docsPlugin), {
    configService: container.configService,
    logger: container.loggerService.logger
  })

  await app.register(controllersPlugin, {
    container
  })

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

  if (process.env.PRODUCTION === 'true') {
    await app.listen({
      host: '0.0.0.0',
      port: container.configService.get('api.port')
    })
  }
}
