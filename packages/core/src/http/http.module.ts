import type { ModuleMetadata } from '#utils/nest'

import { createModule } from '#utils/nest'

export const HttpModuleDefinition: ModuleMetadata = {
  providers: []
}

export const HttpModule = (() =>
  createModule(class HttpModule {}, HttpModuleDefinition))()
