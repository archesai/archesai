import type { FastifyPluginAsync } from 'fastify'

import cors from '@fastify/cors'

// import helmet from '@fastify/helmet'

import type { ConfigService } from '@archesai/core'

export const corsPlugin: FastifyPluginAsync<{
  configService: ConfigService
}> = async (app, { configService }) => {
  const allowedOrigins = configService.get('api.cors.origins').split(',')
  await app.register(cors, {
    allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
    credentials: true,
    maxAge: 86400,
    methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS', 'PATCH'],
    origin: allowedOrigins
  })

  // app.addHook('onSend', async (request, reply) => {
  //   const url = request.url

  //   // API routes - no cache to prevent NS_BINDING_ABORTED
  //   if (url.startsWith('/')) {
  //     reply.header('Cache-Control', 'no-cache, must-revalidate')
  //     reply.header('Pragma', 'no-cache')
  //     reply.header('Expires', '0')
  //   }
  //   // Static assets - long cache
  //   else if (/\.(css|js|png|jpg|jpeg|gif|ico|svg|woff|woff2)$/.exec(url)) {
  //     reply.header('Cache-Control', 'public, max-age=31536000')
  //     reply.header(
  //       'Expires',
  //       new Date(Date.now() + 31536000 * 1000).toUTCString()
  //     )
  //   }
  // })

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
