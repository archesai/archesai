import type { FastifyPluginAsyncZod } from 'fastify-type-provider-zod'

import { ArchesConfigSchema } from '@archesai/schemas'

import type { ConfigService } from '#config/config.service'

export interface ConfigControllerOptions {
  configService: ConfigService
}

export const configController: FastifyPluginAsyncZod<
  ConfigControllerOptions
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
        tags: ['System']
      }
    },
    () => {
      return configService.getConfig()
    }
  )

  await Promise.resolve()
}
