import type { ConfigService } from '@archesai/core'

import type { CreatePortalRequest } from '#portal/dto/create-portal.req.dto'
import type { PortalResource } from '#portal/dto/portal.res.dto'
import type { StripeService } from '#stripe/stripe.service'

/**
 * Service for billing portal.
 */
export class PortalService {
  private readonly configService: ConfigService
  private readonly stripeService: StripeService

  constructor(stripeService: StripeService, configService: ConfigService) {
    this.configService = configService
    this.stripeService = stripeService
  }

  public async create(
    createPortalRequest: CreatePortalRequest
  ): Promise<PortalResource> {
    return this.stripeService.stripe.billingPortal.sessions.create({
      customer: createPortalRequest.organizationId,
      return_url: `${this.configService.get('platform.host')}/organization/billing`
    })
  }
}
