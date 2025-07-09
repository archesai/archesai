import type { EventBus, EventSubscriber } from '@archesai/core'
import type { OrganizationCreatedEvent } from '@archesai/schemas'

import { Logger } from '@archesai/core'

import type { PipelinesService } from '#pipelines/pipelines.service'

/**
 * Subscriber for pipelines.
 */
export class PipelinesSubscriber implements EventSubscriber {
  private readonly eventBus: EventBus
  private readonly logger = new Logger(PipelinesSubscriber.name)
  private readonly pipelinesService: PipelinesService

  constructor(eventBus: EventBus, pipelinesService: PipelinesService) {
    this.eventBus = eventBus
    this.pipelinesService = pipelinesService
  }

  public subscribe() {
    this.eventBus.on(
      'organization.created',
      (event: OrganizationCreatedEvent) => {
        ;(async () => {
          const { organization } = event
          await this.pipelinesService.createDefaultPipeline(
            organization.orgname
          )
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
