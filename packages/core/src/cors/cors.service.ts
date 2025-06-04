import cors from '@fastify/cors'

import type { ConfigService } from '#config/config.service'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

/**
 * Service for setting up CORS.
 */
export class CorsService {
  private readonly configService: ConfigService

  constructor(configService: ConfigService) {
    this.configService = configService
  }

  public setup(app: HttpInstance): void {
    if (this.configService.get('server.cors.enabled')) {
      const allowedOrigins = this.configService
        .get('server.cors.origins')
        .split(',')
      app.register(cors, {
        allowedHeaders: ['Content-Type', 'Authorization'],
        credentials: true,
        maxAge: 86400,
        methods: ['GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE'],
        origin: (origin, cb) => {
          // 1. Same-origin fetches / curl / Postman → origin is undefined → allow
          if (!origin) {
            cb(null, true)
            return
          }

          // 2. Whitelisted sub-domains
          if (allowedOrigins.includes(origin)) {
            cb(null, true)
            return
          }

          // 3. Everything else → block
          cb(new Error('Not allowed by CORS'), false)
        }
      })
    }
  }
}
