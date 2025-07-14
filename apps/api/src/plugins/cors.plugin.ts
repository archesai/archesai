import type { FastifyPluginAsync } from 'fastify'

import cors from '@fastify/cors'

// import helmet from '@fastify/helmet'

import type { ConfigService } from '@archesai/core'

export const corsPlugin: FastifyPluginAsync<{
  configService: ConfigService
}> = async (app, { configService }) => {
  if (configService.get('server.cors.enabled')) {
    const allowedOrigins = configService.get('server.cors.origins').split(',')
    await app.register(cors, {
      allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
      credentials: true,
      maxAge: 86400,
      methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
      origin: allowedOrigins
    })
  }

  // FIXME
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
}
