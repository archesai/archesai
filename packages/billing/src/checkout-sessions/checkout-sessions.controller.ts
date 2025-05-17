import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import { IS_CONTROLLER } from '@archesai/core'

import type { CheckoutSessionsService } from '#checkout-sessions/checkout-sessions.service'
import type { CheckoutSessionResponse } from '#checkout-sessions/dto/checkout-session.res.dto'
import type { CreateCheckoutSessionRequest } from '#checkout-sessions/dto/create-checkout-session.req.dto'

import { CheckoutSessionResponseSchema } from '#checkout-sessions/dto/checkout-session.res.dto'
import { CreateCheckoutSessionRequestSchema } from '#checkout-sessions/dto/create-checkout-session.req.dto'

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
    request: ArchesApiRequest & { body: CreateCheckoutSessionRequest }
  ): Promise<CheckoutSessionResponse> {
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
          body: CreateCheckoutSessionRequestSchema,
          description: 'Create a checkout session',
          operationId: 'createCheckoutSession',
          response: {
            200: CheckoutSessionResponseSchema
          },
          summary: 'Create a checkout session',
          tags: ['Billing - Checkout Sessions']
        }
      },
      this.create.bind(this)
    )
  }
}
