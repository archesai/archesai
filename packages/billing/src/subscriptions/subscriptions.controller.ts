import type { Static } from '@sinclair/typebox'

import { Type } from '@sinclair/typebox'

import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import { IS_CONTROLLER } from '@archesai/core'
import { UpdateSubscriptionDtoSchema } from '@archesai/schemas'

import type { SubscriptionsService } from '#subscriptions/subscriptions.service'

/**
 * Controller for managing subscriptions.
 */
export class SubscriptionsController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly subscriptionsService: SubscriptionsService

  constructor(subscriptionsService: SubscriptionsService) {
    this.subscriptionsService = subscriptionsService
  }

  public async delete(request: ArchesApiRequest & { params: { id: string } }) {
    return this.subscriptionsService.cancel(request.params.id)
  }

  public registerRoutes(app: HttpInstance) {
    app.delete(
      `/billing/subscriptions/:id`,
      {
        schema: {
          description: 'Cancel a subscription',
          operationId: 'cancelSubscription',
          params: Type.Object({
            id: Type.String()
          }),
          summary: 'Cancel a subscription',
          tags: ['Billing']
        }
      },
      this.delete.bind(this)
    )

    app.patch(
      `/billing/subscriptions/:id`,
      {
        schema: {
          body: UpdateSubscriptionDtoSchema,
          description: 'Update a subscription',
          operationId: 'updateSubscription',
          params: Type.Object({
            id: Type.String()
          }),
          summary: 'Update a subscription',
          tags: ['Billing']
        }
      },
      this.update.bind(this)
    )
  }

  public async update(
    request: ArchesApiRequest & {
      body: Static<typeof UpdateSubscriptionDtoSchema>
      params: { id: string }
    }
  ) {
    return this.subscriptionsService.update(
      request.params.id,
      request.body.planId
    )
  }
}
