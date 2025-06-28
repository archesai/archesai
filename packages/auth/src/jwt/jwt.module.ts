import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, createModule } from '@archesai/core'

import { JwtService } from '#jwt/jwt.service'

export const JwtModuleDefinition: ModuleMetadata = {
  exports: [JwtService],
  imports: [ConfigModule],
  providers: [
    {
      inject: [ConfigService],
      provide: JwtService,
      useFactory: (configService: ConfigService) =>
        new JwtService({
          secret: configService.get('jwt.secret')
        })
    }
  ]
}

export const JwtModule = (() =>
  createModule(class JwtModule {}, JwtModuleDefinition))()
