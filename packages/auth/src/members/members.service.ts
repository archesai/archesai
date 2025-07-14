import type { WebsocketsService } from '@archesai/core'
import type { MemberEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { MEMBER_ENTITY_KEY } from '@archesai/schemas'

import type { MemberRepository } from '#members/member.repository'

export const createMembersService = (
  memberRepository: MemberRepository,
  websocketsService: WebsocketsService
) =>
  createBaseService(
    memberRepository,
    websocketsService,
    emitMemberMutationEvent
  )

const emitMemberMutationEvent = (
  entity: MemberEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, MEMBER_ENTITY_KEY]
  })
}

export type MembersService = ReturnType<typeof createMembersService>
