import type { ModuleMetadata } from '@archesai/core'

import {
  DatabaseModule,
  DatabaseService,
  EventBus,
  EventBusModule,
  Module,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'
import { UserEntity } from '@archesai/domain'

import { DeactivatedGuard } from '#users/guards/deactivated.guard'
import { UserRepository } from '#users/user.repository'
import { UsersController } from '#users/users.controller'
import { UsersService } from '#users/users.service'

export const UsersModuleDefinition: ModuleMetadata = {
  exports: [UsersService],
  imports: [DatabaseModule, EventBusModule, WebsocketsModule],
  providers: [
    {
      inject: [EventBus, UserRepository, WebsocketsService],
      provide: UsersService,
      useFactory: (
        eventBus: EventBus,
        userRepository: UserRepository,
        websocketsService: WebsocketsService
      ) => new UsersService(eventBus, userRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: UserRepository,
      useFactory: (databaseService: DatabaseService<UserEntity>) =>
        new UserRepository(databaseService)
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

@Module(UsersModuleDefinition)
export class UsersModule {}
