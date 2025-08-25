import type { BaseService, WebsocketsService } from '@archesai/core'
import type { InvitationEntity, UserEntity } from '@archesai/schemas'

import { createBaseService } from '@archesai/core'
import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import type { InvitationRepository } from '#invitations/invitation.repository'

export const createInvitationsService = (
  invitationRepository: InvitationRepository,
  websocketsService: WebsocketsService
): BaseService<InvitationEntity> & {
  accept(id: string, user: UserEntity): Promise<InvitationEntity>
} => {
  const emitInvitationMutationEvent = (entity: InvitationEntity): void => {
    websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, INVITATION_ENTITY_KEY]
    })
  }
  return {
    ...createBaseService(invitationRepository, emitInvitationMutationEvent),
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

export type InvitationsService = ReturnType<typeof createInvitationsService>
