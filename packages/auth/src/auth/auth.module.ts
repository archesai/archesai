import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { AccessTokensModule } from '#access-tokens/access-tokens.module'
import { AccessTokensService } from '#access-tokens/access-tokens.service'
import { AccountsModule } from '#accounts/accounts.module'
import { AccountsService } from '#accounts/accounts.service'
import { AuthenticationController } from '#auth/auth.controller'
import { AuthenticationService } from '#auth/auth.service'

export const AuthenticationModuleDefinition: ModuleMetadata = {
  imports: [AccessTokensModule, AccountsModule],
  providers: [
    {
      inject: [AccessTokensService],
      provide: AuthenticationService,
      useFactory: (accessTokensService: AccessTokensService) =>
        new AuthenticationService(accessTokensService)
    },
    {
      inject: [AccessTokensService, AccountsService, AuthenticationService],
      provide: AuthenticationController,
      useFactory: (
        accessTokensService: AccessTokensService,
        accountsService: AccountsService,
        authenticationService: AuthenticationService
      ) =>
        new AuthenticationController(
          accessTokensService,
          accountsService,
          authenticationService
        )
    }
  ]
}

@Module(AuthenticationModuleDefinition)
export class AuthenticationModule {}
