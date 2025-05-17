import type { IS_CONTROLLER } from '#common/base.controller'
import type { HttpInstance } from '#http/interfaces/http-instance.interface'

export interface Controller {
  readonly [IS_CONTROLLER]: boolean
  registerRoutes(app: HttpInstance): void
}
