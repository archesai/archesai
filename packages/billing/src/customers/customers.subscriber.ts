import type { EventBus, EventSubscriber } from '@archesai/core'
import type { OrganizationCreatedEvent } from '@archesai/schemas'

import { InternalServerErrorException, Logger } from '@archesai/core'

import type { CustomersService } from '#customers/customers.service'

/**
 * Subscriber for customers.
 */
export class CustomersSubscriber implements EventSubscriber {
  private readonly customersService: CustomersService
  private readonly eventBus: EventBus
  private readonly logger = new Logger(CustomersSubscriber.name)

  constructor(customersService: CustomersService, eventBus: EventBus) {
    this.customersService = customersService
    this.eventBus = eventBus
  }

  public subscribe() {
    this.eventBus.on(
      'organization.created',
      (event: OrganizationCreatedEvent) => {
        ;(async () => {
          this.logger.log('creating customer', {
            orgname: event.organization.orgname
          })

          if (!event.organization.billingEmail) {
            throw new InternalServerErrorException('Billing email is required')
          }

          // replicate your create(...) logic:
          const customer = await this.customersService.create(
            event.organization.orgname,
            event.organization.billingEmail
          )
          this.logger.log('created customer', {
            customer,
            event: event
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
