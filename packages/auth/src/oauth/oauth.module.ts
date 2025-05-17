import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, Module } from '@archesai/core'

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

@Module(OAuthModuleDefinition)
export class OAuthModule {}
