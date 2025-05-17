import type { Controller, HttpInstance } from '@archesai/core'

import { Logger } from '@archesai/core'

/**
 * A loader that registers controllers with the application.
 */
export class ControllerLoader {
  private readonly app: HttpInstance
  private readonly controllers: Controller[]
  private readonly logger = new Logger(ControllerLoader.name)

  constructor(app: HttpInstance, controllers: Controller[]) {
    this.app = app
    this.controllers = controllers
  }

  public loadControllers() {
    this.controllers.forEach((controller) => {
      this.logger.debug(`registering controller ${controller.constructor.name}`)
      controller.registerRoutes(this.app)
    })
  }
}
