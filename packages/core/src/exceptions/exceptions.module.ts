import type { ModuleMetadata } from '#utils/nest'

import { ExceptionsFilter } from '#exceptions/exceptions.filter'
import { Module } from '#utils/nest'

export const ExceptionsModuleDefinition: ModuleMetadata = {
  providers: [
    {
      provide: 'APP_FILTER',
      useClass: ExceptionsFilter
    }
  ]
}

@Module(ExceptionsModuleDefinition)
export class ExceptionsModule {}
