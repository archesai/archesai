import {
  BadRequestException,
  Controller,
  ForbiddenException,
  Param,
  Post,
  Query
} from '@nestjs/common'
import { ApiTags } from '@nestjs/swagger'

import { OrganizationsService } from '@/src/organizations/organizations.service'
import { BillingService } from '@/src/billing/billing.service'
import { BillingUrlEntity } from '@/src/billing/entities/billing-url.entity'
import { ConfigService } from '@/src/config/config.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'
import { RoleTypeEnum } from '@/src/members/entities/member.entity'

@ApiTags('Billing')
@Authenticated([RoleTypeEnum.ADMIN])
@Controller('organizations/:orgname/billing')
export class BillingController {
  constructor(
    private billingService: BillingService,
    private organizationsService: OrganizationsService,
    private configService: ConfigService
  ) {}

  /**
   * Create billing portal
   * @remarks This endpoint will create a billing portal for an organization to edit their subscription and billing information
   * @throws {403} ForbiddenException
   */
  @Post('portal')
  async createBillingPortal(
    @Param('orgname') orgname: string
  ): Promise<BillingUrlEntity> {
    if (!this.configService.get('billing.enabled')) {
      throw new ForbiddenException('Billing is disabled')
    }
    const organization = await this.organizationsService.findOne(orgname)

    return new BillingUrlEntity(
      await this.billingService.createBillingPortal(
        organization.stripeCustomerId
      )
    )
  }

  /**
   * Create checkout session
   * @remarks This endpoint will create a checkout session for an organization to purchase a subscription or one-time product
   * @throws {403} ForbiddenException
   * @throws {400} BadRequestException
   */

  @Post('checkout')
  async createCheckoutSession(
    @Param('orgname') orgname: string,
    @Query('planId') planId: string
  ): Promise<BillingUrlEntity> {
    if (!this.configService.get('billing.enabled')) {
      throw new ForbiddenException('Billing is disabled')
    }
    const organization = await this.organizationsService.findOne(orgname)
    if (['BASIC', 'PREMIUM', 'STANDARD'].includes(organization.plan)) {
      throw new BadRequestException(
        'Cannot purchase a plan when already on a plan'
      )
    }

    const plans = await this.billingService.listPlans()
    const plan = plans.find((p) => p.id === planId)

    if (!plan) {
      throw new BadRequestException('Invalid plan')
    }

    const priceId = plan.priceId

    return new BillingUrlEntity(
      await this.billingService.createCheckoutSession(
        organization.stripeCustomerId,
        {
          price: priceId,
          quantity: 1
        },
        !plan?.recurring?.interval
      )
    )
  }
}
