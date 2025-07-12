import type { WebsocketsService } from '@archesai/core'
import type { InvitationEntity, UserEntity } from '@archesai/schemas'

import { BaseService } from '@archesai/core'
import { INVITATION_ENTITY_KEY } from '@archesai/schemas'

import type { InvitationRepository } from '#invitations/invitation.repository'

/**
 * Service for handling invitations.
 */
export class InvitationsService extends BaseService<InvitationEntity> {
  private readonly invitationRepository: InvitationRepository
  private readonly websocketsService: WebsocketsService

  constructor(
    invitationRepository: InvitationRepository,
    websocketsService: WebsocketsService
  ) {
    super(invitationRepository)
    this.invitationRepository = invitationRepository
    this.websocketsService = websocketsService
  }

  public async accept(id: string, user: UserEntity): Promise<InvitationEntity> {
    const invitation = await this.invitationRepository.findOne(id)
    if (invitation.email !== user.email) {
      throw new Error('Bad Request: Invitation email does not match user email')
    }
    return this.invitationRepository.update(invitation.id, {
      accepted: true
    })
  }

  protected emitMutationEvent(entity: InvitationEntity): void {
    this.websocketsService.broadcastEvent(entity.organizationId, 'update', {
      queryKey: ['organizations', entity.organizationId, INVITATION_ENTITY_KEY]
    })
  }
}
