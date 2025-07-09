import type { WebsocketsService } from '@archesai/core'
import type { MemberEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { MEMBER_ENTITY_KEY } from '@archesai/schemas'

import type { MemberRepository } from '#members/member.repository'

/**
 * Service for handling members.
 */
export class MembersService extends BaseService<MemberEntity> {
  private readonly websocketsService: WebsocketsService
  constructor(
    memberRepository: MemberRepository,
    websocketsService: WebsocketsService
  ) {
    super(memberRepository)
    this.websocketsService = websocketsService
  }

  protected emitMutationEvent(entity: MemberEntity): void {
    this.websocketsService.broadcastEvent(entity.orgname, 'update', {
      queryKey: ['organizations', entity.orgname, MEMBER_ENTITY_KEY]
    })
  }
}
