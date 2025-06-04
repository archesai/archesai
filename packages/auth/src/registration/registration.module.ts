import type { ModuleMetadata } from '@archesai/core'

import { Module } from '@archesai/core'

import { AccountsModule } from '#accounts/accounts.module'
import { AccountsService } from '#accounts/accounts.service'
import { RegistrationController } from '#registration/registration.controller'
import { RegistrationService } from '#registration/registration.service'
import { UsersModule } from '#users/users.module'
import { UsersService } from '#users/users.service'

export const RegistrationModuleDefinition: ModuleMetadata = {
  imports: [AccountsModule, UsersModule],
  providers: [
    {
      inject: [AccountsService, UsersService],
      provide: RegistrationService,
      useFactory: (
        accountsService: AccountsService,
        usersService: UsersService
      ) => new RegistrationService(accountsService, usersService)
    },
    {
      inject: [RegistrationService],
      provide: RegistrationController,
      useFactory: (registrationService: RegistrationService) =>
        new RegistrationController(registrationService)
    }
  ]
}

@Module(RegistrationModuleDefinition)
export class RegistrationModule {}
