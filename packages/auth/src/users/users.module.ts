import type { ModuleMetadata } from '@archesai/core'
import type { UserInsertModel, UserSelectModel } from '@archesai/database'
import type { UserEntity } from '@archesai/schemas'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { OrganizationsModule } from '#organizations/organizations.module'
import { OrganizationsService } from '#organizations/organizations.service'
import { DeactivatedGuard } from '#users/guards/deactivated.guard'
import { UserRepository } from '#users/user.repository'
import { UsersController } from '#users/users.controller'
import { UsersService } from '#users/users.service'

export const UsersModuleDefinition: ModuleMetadata = {
  exports: [UsersService],
  imports: [DatabaseModule, OrganizationsModule, WebsocketsModule],
  providers: [
    {
      inject: [OrganizationsService, UserRepository, WebsocketsService],
      provide: UsersService,
      useFactory: (
        organizationsService: OrganizationsService,
        userRepository: UserRepository,
        websocketsService: WebsocketsService
      ) =>
        new UsersService(
          organizationsService,
          userRepository,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: UserRepository,
      useFactory: (
        databaseService: DatabaseService<
          UserEntity,
          UserInsertModel,
          UserSelectModel
        >
      ) => new UserRepository(databaseService)
    },
    {
      provide: DeactivatedGuard,
      useClass: DeactivatedGuard
    },
    {
      inject: [UsersService],
      provide: UsersController,
      useFactory: (usersService: UsersService) =>
        new UsersController(usersService)
    }
  ]
}

export const UsersModule = (() =>
  createModule(class UsersModule {}, UsersModuleDefinition))()
