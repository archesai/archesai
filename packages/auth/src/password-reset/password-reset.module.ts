import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { AccountsModule } from '#accounts/accounts.module'
import { AccountsService } from '#accounts/accounts.service'
import { HashingModule } from '#hashing/hashing.module'
import { HashingService } from '#hashing/hashing.service'
import { PasswordResetController } from '#password-reset/password-reset.controller'
import { PasswordResetService } from '#password-reset/password-reset.service'
import { VerificationTokensModule } from '#verification-tokens/verification-tokens.module'
import { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

export const PasswordResetModuleDefinition: ModuleMetadata = {
  imports: [AccountsModule, HashingModule, VerificationTokensModule],
  providers: [
    {
      inject: [AccountsService, HashingService, VerificationTokensService],
      provide: PasswordResetService,
      useFactory: (
        accountsService: AccountsService,
        hashingService: HashingService,
        verificationTokensService: VerificationTokensService
      ) =>
        new PasswordResetService(
          accountsService,
          hashingService,
          verificationTokensService
        )
    },
    {
      inject: [PasswordResetService],
      provide: PasswordResetController,
      useFactory: (passwordResetService: PasswordResetService) =>
        new PasswordResetController(passwordResetService)
    }
  ]
}

@Module(PasswordResetModuleDefinition)
export class PasswordResetModule {}
