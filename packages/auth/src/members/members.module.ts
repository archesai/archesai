import type { ModuleMetadata } from '@archesai/core'
import type { MemberEntity } from '@archesai/schemas'

import {
  createModule,
  DatabaseModule,
  DatabaseService,
  WebsocketsModule,
  WebsocketsService
} from '@archesai/core'

import { MembershipGuard } from '#members/guards/membership.guard'
import { MemberRepository } from '#members/member.repository'
import { MembersController } from '#members/members.controller'
import { MembersService } from '#members/members.service'

export const MembersModuleDefinition: ModuleMetadata = {
  exports: [MembersService],
  imports: [DatabaseModule, WebsocketsModule],
  providers: [
    {
      inject: [MemberRepository, WebsocketsService],
      provide: MembersService,
      useFactory: (
        memberRepository: MemberRepository,
        websocketsService: WebsocketsService
      ) => new MembersService(memberRepository, websocketsService)
    },
    {
      inject: [DatabaseService],
      provide: MemberRepository,
      useFactory: (databaseService: DatabaseService<MemberEntity>) =>
        new MemberRepository(databaseService)
    },
    {
      inject: [MembersService],
      provide: MembershipGuard,
      useFactory: (membersService: MembersService) => {
        return new MembershipGuard(membersService)
      }
    },
    {
      inject: [MembersService],
      provide: MembersController,
      useFactory: (membersService: MembersService) =>
        new MembersController(membersService)
    }
  ]
}

export const MembersModule = (() =>
  createModule(class MembersModule {}, MembersModuleDefinition))()
