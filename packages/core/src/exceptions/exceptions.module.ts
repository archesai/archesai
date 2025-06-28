import type { ModuleMetadata } from '#utils/nest'

import { ExceptionsFilter } from '#exceptions/exceptions.filter'
import { createModule } from '#utils/nest'

export const ExceptionsModuleDefinition: ModuleMetadata = {
  providers: [
    {
      provide: 'APP_FILTER',
      useClass: ExceptionsFilter
    }
  ]
}

export const ExceptionsModule = (() =>
  createModule(class ExceptionsModule {}, ExceptionsModuleDefinition))()
