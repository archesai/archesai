import type { ModuleMetadata } from '@archesai/core'

import { createModule } from '@archesai/core'

import { EmailChangeController } from '#email-change/email-change.controller'
import { EmailChangeService } from '#email-change/email-change.service'
import { UsersModule } from '#users/users.module'
import { UsersService } from '#users/users.service'
import { VerificationTokensModule } from '#verification-tokens/verification-tokens.module'
import { VerificationTokensService } from '#verification-tokens/verification-tokens.service'

export const EmailChangeModuleDefinition: ModuleMetadata = {
  imports: [UsersModule, VerificationTokensModule],
  providers: [
    {
      inject: [EmailChangeService],
      provide: EmailChangeController,
      useFactory: (emailChangeService: EmailChangeService) =>
        new EmailChangeController(emailChangeService)
    },
    {
      inject: [UsersService, VerificationTokensService],
      provide: EmailChangeService,
      useFactory: (
        usersService: UsersService,
        verificationTokensService: VerificationTokensService
      ) => new EmailChangeService(usersService, verificationTokensService)
    }
  ]
}

export const EmailChangeModule = (() =>
  createModule(class EmailChangeModule {}, EmailChangeModuleDefinition))()
