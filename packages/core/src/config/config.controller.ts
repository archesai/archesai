import type { ConfigService } from '#config/config.service'
import type { Controller } from '#http/interfaces/controller.interface'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

import { IS_CONTROLLER } from '#common/base.controller'
import { ArchesConfigSchema } from '#config/schemas/config.schema'

/**
 * Controller for managing the application configuration.
 */
export class ConfigController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly configService: ConfigService

  constructor(configService: ConfigService) {
    this.configService = configService
  }

  public config() {
    return this.configService.getConfig()
  }

  public registerRoutes(app: HttpInstance) {
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
      this.config.bind(this)
    )
  }
}
