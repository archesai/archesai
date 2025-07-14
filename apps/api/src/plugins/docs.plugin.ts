import type { SwaggerOptions } from '@fastify/swagger'
import type { FastifyPluginAsync } from 'fastify'

import fastifySwagger from '@fastify/swagger'
import scalarUi from '@scalar/fastify-api-reference'

import type { ConfigService, Logger } from '@archesai/core'

import {
  ErrorDocumentSchema,
  ErrorObjectSchema,
  FieldFilterSchema
} from '@archesai/core'

export const docsPlugin: FastifyPluginAsync<{
  configService: ConfigService
  logger: Logger
}> = async (app, { configService, logger }) => {
  const DEFAULT_MODELS = [
    FieldFilterSchema,
    ErrorObjectSchema,
    ErrorDocumentSchema
  ]
  if (configService.get('server.docs.enabled')) {
    // Register Default Schemas
    for (const model of DEFAULT_MODELS) {
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
            url: `${configService.get('tls.enabled') ? 'https://' : 'http://'}${configService.get('server.host')}`
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
        forceDarkModeState: 'dark',
        layout: 'modern',
        theme: 'purple',
        title: 'Arches AI API',
        withDefaultFonts: true
      },
      routePrefix: '/docs'
    })
  }
  logger.debug('documentation setup complete')
}
