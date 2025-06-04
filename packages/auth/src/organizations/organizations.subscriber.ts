import type { EventBus, EventSubscriber } from '@archesai/core'
import type { UserCreatedEvent } from '@archesai/domain'

import { Logger } from '@archesai/core'

import type { OrganizationsService } from '#organizations/organizations.service'

/**
 * Subscriber for handling organizations.
 */
export class OrganizationsSubscriber implements EventSubscriber {
  private readonly eventBus: EventBus
  private readonly logger = new Logger(OrganizationsSubscriber.name)
  private readonly organizationsService: OrganizationsService

  constructor(eventBus: EventBus, organizationsService: OrganizationsService) {
    this.eventBus = eventBus
    this.organizationsService = organizationsService
  }

  public subscribe() {
    this.eventBus.on('user.created', (event: UserCreatedEvent) => {
      ;(async () => {
        const { user } = event
        const usernamePrefix = user.email.split('@')[0]
        if (!usernamePrefix) {
          throw new Error('Could not extract username from email')
        }
        const orgname = usernamePrefix + user.id.slice(0, 5)
        await this.organizationsService.create({
          billingEmail: user.email,
          credits: 0,
          orgname,
          plan: 'FREE'
        })
        this.logger.log('user created event', { user })
      })().catch((error: unknown) => {
        this.logger.error('error', {
          error,
          event
        })
      })
    })
  }
}
