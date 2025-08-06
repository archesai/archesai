import type { SwaggerOptions } from '@fastify/swagger'
import type { FastifyPluginAsync } from 'fastify'

import fastifySwagger from '@fastify/swagger'
import scalarUi from '@scalar/fastify-api-reference'
import {
  jsonSchemaTransform,
  jsonSchemaTransformObject
} from 'fastify-type-provider-zod'

import type { ConfigService, Logger } from '@archesai/core'

export const docsPlugin: FastifyPluginAsync<{
  configService: ConfigService
  logger: Logger
}> = async (app, { configService, logger }) => {
  if (configService.get('api.docs')) {
    // Register fastify plugin
    await app.register(fastifySwagger, {
      openapi: {
        components: {
          securitySchemes: {
            bearerAuth: {
              bearerFormat: 'JWT',
              description: 'API Token for authenticated requests',
              scheme: 'bearer',
              type: 'http'
            },
            sessionCookie: {
              description: 'Session cookie for authenticated requests',
              in: 'cookie',
              name: '__Secure-better-auth.session_token',
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
            url: `http://${configService.get('api.host')}`
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
      },
      transform: jsonSchemaTransform,
      transformObject: jsonSchemaTransformObject
    } satisfies SwaggerOptions)

    // Register scalar-ui plugin
    await app.register(scalarUi, {
      configuration: {
        _integration: 'fastify',
        forceDarkModeState: 'dark',
        hideModels: true,
        layout: 'modern',
        pageTitle: 'Arches AI API',
        persistAuth: true,
        tagsSorter: 'alpha',
        theme: 'purple',
        title: 'Arches AI API',
        withDefaultFonts: true
      },
      routePrefix: '/docs'
    })
  }
  logger.debug('documentation setup complete')
}
