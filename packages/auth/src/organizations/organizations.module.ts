import type { ModuleMetadata } from '@archesai/core'
import type { OrganizationEntity } from '@archesai/domain'

import {
  ConfigModule,
  ConfigService,
  createModule,
  DatabaseModule,
  DatabaseService,
  EventBus,
  EventBusModule,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { OrganizationRepository } from '#organizations/organization.repository'
import { OrganizationsController } from '#organizations/organizations.controller'
import { OrganizationsService } from '#organizations/organizations.service'
import { OrganizationsSubscriber } from '#organizations/organizations.subscriber'

export const OrganizationsModuleDefinition: ModuleMetadata = {
  exports: [OrganizationsService],
  imports: [ConfigModule, DatabaseModule, EventBusModule, WebsocketsModule],
  providers: [
    {
      inject: [
        ConfigService,
        EventBus,
        OrganizationRepository,
        WebsocketsService
      ],
      provide: OrganizationsService,
      useFactory: (
        configService: ConfigService,
        eventBus: EventBus,
        organizationRepository: OrganizationRepository,
        websocketsService: WebsocketsService
      ) =>
        new OrganizationsService(
          configService,
          eventBus,
          organizationRepository,
          websocketsService
        )
    },
    {
      inject: [DatabaseService],
      provide: OrganizationRepository,
      useFactory: (databaseService: DatabaseService<OrganizationEntity>) =>
        new OrganizationRepository(databaseService)
    },
    {
      inject: [EventBus, OrganizationsService],
      provide: OrganizationsSubscriber,
      useFactory: (
        eventBus: EventBus,
        organizationsService: OrganizationsService
      ) => new OrganizationsSubscriber(eventBus, organizationsService)
    },
    {
      inject: [OrganizationsService],
      provide: OrganizationsController,
      useFactory: (organizationsService: OrganizationsService) =>
        new OrganizationsController(organizationsService)
    }
  ]
}

export const OrganizationsModule = (() =>
  createModule(class OrganizationsModule {}, OrganizationsModuleDefinition))()
