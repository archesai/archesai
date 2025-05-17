import type { SwaggerOptions } from '@fastify/swagger'

import fastifySwagger from '@fastify/swagger'
import scalarUi from '@scalar/fastify-api-reference'

import type { ConfigService } from '#config/config.service'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

import { ArchesApiForbiddenResponseSchema } from '#exceptions/schemas/arches-api-forbidden-response.schema'
import { ArchesApiNoContentResponseSchema } from '#exceptions/schemas/arches-api-no-content-response.schema'
import { ArchesApiNotFoundResponseSchema } from '#exceptions/schemas/arches-api-not-found-response.schema'
import { ArchesApiUnauthorizedResponseSchema } from '#exceptions/schemas/arches-api-unauthorized-response.schema'
import { FieldFilterSchema } from '#http/dto/search-query.dto'
import { ApiResponseSchema } from '#http/schemas/api-response.schema'
import { ErrorsSchema } from '#http/schemas/errors.schema'
import { IncludedSchema } from '#http/schemas/included.schema'
import { MetaSchema } from '#http/schemas/meta.schema'
import { Logger } from '#logging/logger'

/**
 * Service for setting up the API documentation.
 */
export class DocsService {
  private readonly configService: ConfigService
  private readonly DEFAULT_MODELS = [
    FieldFilterSchema,
    ApiResponseSchema,
    ErrorsSchema,
    IncludedSchema,
    MetaSchema,
    ArchesApiForbiddenResponseSchema,
    ArchesApiNoContentResponseSchema,
    ArchesApiUnauthorizedResponseSchema,
    ArchesApiNotFoundResponseSchema
  ]
  private readonly logger = new Logger(DocsService.name)
  constructor(configService: ConfigService) {
    this.configService = configService
  }

  public async setup(app: HttpInstance): Promise<void> {
    this.logger.debug('setting up documentation')
    if (this.configService.get('server.docs.enabled')) {
      // Register Default Schemas
      for (const model of this.DEFAULT_MODELS) {
        app.addSchema(model)
      }

      // Register fastify plugin
      await app.register(fastifySwagger, {
        openapi: {
          components: {
            securitySchemes: {
              bearerAuth: {
                bearerFormat: 'JWT',
                description: 'API Token',
                scheme: 'bearer',
                type: 'http'
              },
              sessionCookie: {
                description: 'Session Cookie',
                in: 'cookie',
                name: 'sessionid',
                type: 'apiKey'
              }
            }
          },
          info: {
            description: 'The Arches AI API',
            title: 'Arches AI API',
            version: 'v1'
          },
          jsonSchemaDialect: 'https://spec.openapis.org/oas/3.1/dialect/base',
          openapi: '3.1.1',
          servers: [
            {
              url: `${this.configService.get('tls.enabled') ? 'https://' : 'http://'}${this.configService.get('server.host')}:${this.configService.get('server.port').toString()}`
            }
          ]
        },
        refResolver: {
          buildLocalReference(json, _baseUri, _fragment, i) {
            if (!json.title && json.$id) {
              json.title = json.$id
            }
            // Fallback if no $id is present
            if (!json.$id || typeof json.$id !== 'string') {
              return `def-${i.toString()}`
            }

            return json.$id
          }
        }
      } satisfies SwaggerOptions)

      // Register scalar-ui plugin
      await app.register(scalarUi, {
        configuration: {
          _integration: 'fastify',
          baseServerURL: `${this.configService.get('tls.enabled') ? 'https://' : 'http://'}${this.configService.get('server.host')}:${this.configService.get('server.port').toString()}`,
          forceDarkModeState: 'dark',
          layout: 'modern',
          theme: 'purple',
          title: 'Arches AI API',
          withDefaultFonts: true
        },
        routePrefix: '/docs'
      })
    }
    this.logger.debug('documentation setup complete')
  }
}
