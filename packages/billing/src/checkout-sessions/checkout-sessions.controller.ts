import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'
import type {
  CheckoutSessionDto,
  CreateCheckoutSessionDto
} from '@archesai/schemas'

import { IS_CONTROLLER } from '@archesai/core'
import {
  CheckoutSessionDtoSchema,
  CreateCheckoutSessionDtoSchema
} from '@archesai/schemas'

import type { CheckoutSessionsService } from '#checkout-sessions/checkout-sessions.service'

/**
 * Controller for checkout sessions.
 */
export class CheckoutSessionsController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly checkoutSessionsService: CheckoutSessionsService

  constructor(checkoutSessionsService: CheckoutSessionsService) {
    this.checkoutSessionsService = checkoutSessionsService
  }

  public async create(
    request: ArchesApiRequest & { body: CreateCheckoutSessionDto }
  ): Promise<CheckoutSessionDto> {
    return this.checkoutSessionsService.create(
      request.user!.orgname,
      {
        price: request.body.priceId,
        quantity: 1
      },
      false
    )
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/billing/checkout-sessions`,
      {
        schema: {
          body: CreateCheckoutSessionDtoSchema,
          description: 'Create a checkout session',
          operationId: 'createCheckoutSession',
          response: {
            200: CheckoutSessionDtoSchema
          },
          summary: 'Create a checkout session',
          tags: ['Billing']
        }
      },
      this.create.bind(this)
    )
  }
}
