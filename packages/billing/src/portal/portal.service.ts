import type { ConfigService } from '@archesai/core'
import type { CreatePortalDto, PortalDto } from '@archesai/schemas'

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
    createPortalRequest: CreatePortalDto
  ): Promise<PortalDto> {
    return this.stripeService.stripe.billingPortal.sessions.create({
      customer: createPortalRequest.organizationId,
      return_url: `${this.configService.get('platform.host')}/organization/billing`
    })
  }
}
