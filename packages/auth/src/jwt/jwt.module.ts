import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, Module } from '@archesai/core'

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

@Module(JwtModuleDefinition)
export class JwtModule {}
