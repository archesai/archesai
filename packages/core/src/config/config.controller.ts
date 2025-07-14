import type { FastifyPluginAsyncTypebox } from '@fastify/type-provider-typebox'

import type { ConfigService } from '#config/config.service'

import { ArchesConfigSchema } from '#config/schemas/config.schema'

export interface ConfigControllerOptions {
  configService: ConfigService
}

export const configController: FastifyPluginAsyncTypebox<
  ConfigControllerOptions
  // eslint-disable-next-line @typescript-eslint/require-await
> = async (app, { configService }) => {
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
}
