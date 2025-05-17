import type { ArchesApiRequest, Controller, HttpInstance } from '@archesai/core'

import { InternalServerErrorException, IS_CONTROLLER } from '@archesai/core'

import type { CallbacksService } from '#callbacks/callbacks.service'

/**
 * Controller for callbacks.
 */
export class CallbacksController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly callbacksService: CallbacksService

  constructor(callbacksService: CallbacksService) {
    this.callbacksService = callbacksService
  }

  public async callback(request: ArchesApiRequest) {
    const stripeSignature = request.headers['stripe-signature']
    if (typeof stripeSignature !== 'string') {
      throw new InternalServerErrorException('Missing stripe-signature header')
    }
    await this.callbacksService.handle(stripeSignature, request)
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/billing/stripe/callback`,
      {
        schema: {
          description: `Handles Stripe webhook callbacks`,
          operationId: 'stripeCallback',
          response: {
            200: {
              description: 'OK'
            }
          },
          summary: `Stripe webhook callback`,
          tags: ['Billing - Callbacks']
        }
      },
      this.callback.bind(this)
    )
  }
}
