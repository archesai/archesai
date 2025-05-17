import type { EventBus, EventSubscriber } from '@archesai/core'
import type { OrganizationCreatedEvent } from '@archesai/domain'

import { Logger } from '@archesai/core'

import type { InvitationsService } from '#invitations/invitations.service'
import type { MembersService } from '#members/members.service'
import type { UsersService } from '#users/users.service'

/**
 * Subscribes to organization created events and creates an accepted invitation and membership for the creator.
 */
export class InvitationsSubscriber implements EventSubscriber {
  private readonly eventBus: EventBus
  private readonly invitationsService: InvitationsService
  private readonly logger = new Logger(InvitationsSubscriber.name)
  private readonly membersService: MembersService
  private readonly usersService: UsersService

  constructor(
    eventBus: EventBus,
    invitationsService: InvitationsService,
    membersService: MembersService,
    usersService: UsersService
  ) {
    this.eventBus = eventBus
    this.invitationsService = invitationsService
    this.membersService = membersService
    this.usersService = usersService
  }

  public subscribe() {
    this.eventBus.on(
      'organization.created',
      (event: OrganizationCreatedEvent) => {
        ;(async () => {
          if (!event.creator) {
            this.logger.warn(`no creator found`)
            return
          }

          const user = await this.usersService.findOne(event.creator.id)
          const invitation = await this.invitationsService.create({
            accepted: true,
            email: user.email,
            orgname: event.organization.orgname,
            role: 'ADMIN'
          })
          const membership = await this.membersService.create({
            invitationId: invitation.id,
            name: event.creator.name,
            orgname: event.organization.orgname,
            role: 'USER',
            userId: event.creator.id
          })

          this.logger.log(`created membership`, {
            membership
          })
        })().catch((error: unknown) => {
          this.logger.error('error', {
            error,
            event
          })
        })
      }
    )
  }
}
