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
            callback(new Error('Not allowed by CORS'), false)
          }
        }
      })
    }
  }
}
