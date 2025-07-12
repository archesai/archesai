import type { Controller, HttpInstance } from '@archesai/core'

import { IS_CONTROLLER } from '@archesai/core'
import {
  CheckoutSessionDtoSchema,
  CreateCheckoutSessionDtoSchema,
  CreatePortalDtoSchema,
  PortalDtoSchema
} from '@archesai/schemas'

import type { StripeService } from '#stripe/stripe.service'

/**
 * Controller for billing portal.
 */
export class StripeController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService) {
    this.stripeService = stripeService
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/billing/portal`,
      {
        schema: {
          body: CreatePortalDtoSchema,
          description: 'Create a new portal',
          operationId: 'createPortal',
          response: {
            201: {
              description: 'The created portal',
              schema: PortalDtoSchema
            }
          },
          summary: 'Create a new portal',
          tags: ['Billing']
        }
      },
      (req) => {
        return this.stripeService.createPortal(req.body)
      }
    )

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
      (req) => {
        return this.stripeService.createCheckoutSession(
          '', // FIXME should be req.user!.organizationId,
          {
            price: req.body.priceId,
            quantity: 1
          },
          false
        )
      }
    )
  }
}
