import type { ModuleMetadata } from '@archesai/core'

import { createModule } from '@archesai/core'

import { EmailVerificationController } from '#email-verification/email-verification.controller'
import { EmailVerificationService } from '#email-verification/email-verification.service'
import { UsersModule } from '#users/users.module'
import { UsersService } from '#users/users.service'
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

export const EmailVerificationModule = (() =>
  createModule(
    class EmailVerificationModule {},
    EmailVerificationModuleDefinition
  ))()
