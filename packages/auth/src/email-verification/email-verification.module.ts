import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'
import { UsersModule, UsersService } from '@archesai/organizations'

import { EmailVerificationController } from '#email-verification/email-verification.controller'
import { EmailVerificationService } from '#email-verification/email-verification.service'
import { VerificationTokensModule } from '#verification-tokens/verification-tokens.module'
import { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

export const EmailVerificationModuleDefinition: ModuleMetadata = {
  imports: [UsersModule, VerificationTokensModule],
  providers: [
    {
      inject: [UsersService, VerificationTokensService],
      provide: EmailVerificationService,
      useFactory: (
        usersService: UsersService,
        verificationTokensService: VerificationTokensService
      ) => new EmailVerificationService(usersService, verificationTokensService)
    },
    {
      inject: [EmailVerificationService],
      provide: EmailVerificationController,
      useFactory: (emailVerificationService: EmailVerificationService) =>
        new EmailVerificationController(emailVerificationService)
    }
  ]
}

@Module(EmailVerificationModuleDefinition)
export class EmailVerificationModule {}
