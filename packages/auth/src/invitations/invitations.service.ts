import type { WebsocketsService } from '@archesai/core'
import type { InvitationEntity, UserEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import type { InvitationRepository } from '#invitations/invitation.repository'

export const createInvitationsService = (
  invitationRepository: InvitationRepository,
  websocketsService: WebsocketsService
) => {
  return {
    ...createBaseService(
      invitationRepository,
      websocketsService,
      emitInvitationMutationEvent
    ),
    async accept(id: string, user: UserEntity): Promise<InvitationEntity> {
      const invitation = await invitationRepository.findOne(id)
      if (invitation.email !== user.email) {
        throw new Error(
          'Bad Request: Invitation email does not match user email'
        )
      }
      return invitationRepository.update(invitation.id, {
        status: 'accepted'
      })
    }
  }
}

const emitInvitationMutationEvent = (
  entity: InvitationEntity,
  websocketsService: WebsocketsService
): void => {
  websocketsService.broadcastEvent(entity.organizationId, 'update', {
    queryKey: ['organizations', entity.organizationId, INVITATION_ENTITY_KEY]
  })
}

export type InvitationsService = ReturnType<typeof createInvitationsService>
