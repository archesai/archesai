import type { NestFastifyApplication } from '@nestjs/platform-fastify'

// import helmet from '@fastify/helmet'
import { DiscoveryService, NestFactory } from '@nestjs/core'
import { FastifyAdapter } from '@nestjs/platform-fastify'

import type { Controller, HttpInstance } from '@archesai/core'

import { SessionsService } from '@archesai/auth'
import {
  ControllerLoader,
  CorsService,
  DocsService,
  Logger,
  WebsocketsService
} from '@archesai/core'

import { AppModule } from '#app/app.module'

export async function setup(): Promise<NestFastifyApplication> {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule.forRoot(),
    new FastifyAdapter({
      disableRequestLogging: true
    })
  )

  // Get fastify instance
  const httpInstance = app
    .getHttpAdapter()
    .getInstance() as unknown as HttpInstance

  //  Setup Logger
  app.useLogger(new Logger('NestApplication'))

  // Docs Setup
  const docsService = app.get(DocsService)
  await docsService.setup(httpInstance)

  // CORS Configuration
  const corsService = app.get(CorsService)
  corsService.setup(httpInstance)

  // // Security Middlewares
  // await httpInstance.register(helmet, {
  //   contentSecurityPolicy: {
  //     directives: {
  //       defaultSrc: [`'self'`],
  //       fontSrc: [`'self'`, 'fonts.scalar.com', 'data:'],
  //       imgSrc: [`'self'`, 'data:'],
  //       scriptSrc: [`'self'`, `https: 'unsafe-inline'`, `'unsafe-eval'`],
  //       styleSrc: [`'self'`, `'unsafe-inline'`, 'fonts.scalar.com']
  //     }
  //   }
  // })

  // Session Management
  const sessionsService = app.get(SessionsService)
  sessionsService.setup(httpInstance)

  // Websocket Adapter
  const websocketsService = app.get(WebsocketsService)
  await websocketsService.setupWebsocketAdapter(app.getHttpServer())

  // Discover all controllers automatically
  const discoveryService = app.get(DiscoveryService)
  const controllers = discoveryService
    .getProviders()
    .map((wrapper) => wrapper.instance as unknown)
    .filter(
      (instance): instance is Controller =>
        instance !== null &&
        typeof instance === 'object' &&
        Symbol.for('isController') in instance
    )

  // Instantiate and load routes with your custom ControllerLoader
  const controllerLoader = new ControllerLoader(httpInstance, controllers)
  controllerLoader.loadControllers()

  // Enable Shutdown Hooks
  app.enableShutdownHooks()

  return app
}
