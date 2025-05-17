import type { ModuleMetadata } from '#utils/nest'

import { Module } from '#utils/nest'

export const HttpModuleDefinition: ModuleMetadata = {
  providers: []
}

@Module(HttpModuleDefinition)
export class HttpModule {}
