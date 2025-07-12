import type { ModuleMetadata } from '@archesai/core'
import type {
  InvitationInsertModel,
  InvitationSelectModel
} from '@archesai/database'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  EventBus,
  EventBusModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { InvitationRepository } from '#invitations/invitation.repository'
import { InvitationsController } from '#invitations/invitations.controller'
import { InvitationsService } from '#invitations/invitations.service'
import { InvitationsSubscriber } from '#invitations/invitations.subscriber'
import { MembersModule } from '#members/members.module'
import { MembersService } from '#members/members.service'
import { UsersModule } from '#users/users.module'
import { UsersService } from '#users/users.service'

export const InvitationsModuleDefinition: ModuleMetadata = {
  exports: [InvitationsService],
  imports: [
    DatabaseModule,
    EventBusModule,
    MembersModule,
    UsersModule,
    WebsocketsModule
  ],
  providers: [
    {
      inject: [InvitationRepository, WebsocketsService],
      provide: InvitationsService,
      useFactory: (
        invitationRepository: InvitationRepository,
        websocketsService: WebsocketsService
      ) => new InvitationsService(invitationRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: InvitationRepository,
      useFactory: (
        databaseService: DatabaseService<
          InvitationInsertModel,
          InvitationSelectModel
        >
      ) => new InvitationRepository(databaseService)
    },
    {
      inject: [EventBus, InvitationsService, MembersService, UsersService],
      provide: InvitationsSubscriber,
      useFactory: (
        eventBus: EventBus,
        invitationsService: InvitationsService,
        membersService: MembersService,
        usersService: UsersService
      ) =>
        new InvitationsSubscriber(
          eventBus,
          invitationsService,
          membersService,
          usersService
        )
    },
    {
      inject: [InvitationsService],
      provide: InvitationsController,
      useFactory: (invitationsService: InvitationsService) =>
        new InvitationsController(invitationsService)
    }
  ]
}

export const InvitationsModule = (() =>
  createModule(class InvitationsModule {}, InvitationsModuleDefinition))()
