import type { ModuleMetadata } from '@archesai/core'

import { ConfigModule, ConfigService, Module } from '@archesai/core'

import { AccessTokensService } from '#access-tokens/access-tokens.service'
import { AccountsModule } from '#accounts/accounts.module'
import { AccountsService } from '#accounts/accounts.service'
import { JwtModule } from '#jwt/jwt.module'
import { JwtService } from '#jwt/jwt.service'

export const AccessTokensModuleDefinition: ModuleMetadata = {
  exports: [AccessTokensService],
  imports: [AccountsModule, ConfigModule, JwtModule],
  providers: [
    {
      inject: [AccountsService, ConfigService, JwtService],
      provide: AccessTokensService,
      useFactory: (accountsService: AccountsService, jwtService: JwtService) =>
        new AccessTokensService(accountsService, jwtService)
    }
  ]
}

@Module(AccessTokensModuleDefinition)
export class AccessTokensModule {}
