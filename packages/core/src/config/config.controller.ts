import type { FastifyPluginCallbackTypebox } from '@fastify/type-provider-typebox'

import type { ConfigService } from '#config/config.service'

import { ArchesConfigSchema } from '#config/schemas/config.schema'

export interface ConfigControllerOptions {
  configService: ConfigService
}

export const configController: FastifyPluginCallbackTypebox<
  ConfigControllerOptions
> = (app, { configService }, done) => {
  app.get(
    `/config`,
    {
      schema: {
        description: `Get the configuration of the application`,
        operationId: 'getConfig',
        response: {
          200: ArchesConfigSchema
        },
        summary: `Get the configuration`,
        tags: ['Configuration']
      }
    },
    () => {
      return configService.getConfig()
    }
  )

  done()
}
