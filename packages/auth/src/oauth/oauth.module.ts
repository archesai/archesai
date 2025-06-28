import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, createModule } from '@archesai/core'

import { OAuthController } from '#oauth/oauth.controller'
import { OAuthService } from '#oauth/oauth.service'

export const OAuthModuleDefinition: ModuleMetadata = {
  imports: [ConfigModule],
  providers: [
    {
      provide: OAuthService,
      useFactory: () => new OAuthService()
    },
    {
      provide: OAuthController,
      useFactory: () => new OAuthController()
    }
  ]
}

export const OAuthModule = (() =>
  createModule(class OAuthModule {}, OAuthModuleDefinition))()
