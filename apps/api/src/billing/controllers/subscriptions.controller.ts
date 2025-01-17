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
import { ConfigService } from '@/src/config/config.service'
import { Authenticated } from '@/src/auth/decorators/authenticated.decorator'
import { RoleTypeEnum } from '@/src/members/entities/member.entity'

@ApiTags('Billing - Subscription')
@Authenticated([RoleTypeEnum.ADMIN])
@Controller('/organizations/:orgname/billing/subscription')
export class SubscriptionsController {
  constructor(
    private billingService: BillingService,
    private organizationsService: OrganizationsService,
    private configService: ConfigService
  ) {}

  /**
   * Change subscription plan
   * @remarks This endpoint will change the subscription plan for an organization
   * @throws {403} ForbiddenException
   * @throws {400} BadRequestException
   */
  @Post()
  async update(
    @Param('orgname') orgname: string,
    @Query('planId') planId: string
  ) {
    if (!this.configService.get('billing.enabled')) {
      throw new ForbiddenException('Billing is disabled')
    }
    const organization = await this.organizationsService.findOne(orgname)

    const plans = await this.billingService.listPlans()
    const plan = plans.find((p) => p.id === planId)

    if (!plan) {
      throw new BadRequestException('Invalid plan')
    }

    // Update the subscription to the new plan
    await this.billingService.updateSubscription(
      organization.stripeCustomerId,
      plan.priceId
    )
  }

  /**
   * Cancel subscription plan
   * @remarks This endpoint will cancel the subscription plan for an organization
   * @throws {403} ForbiddenException
   */
  @Post('cancel')
  async cancelSubscriptionPlan(@Param('orgname') orgname: string) {
    if (!this.configService.get('billing.enabled')) {
      throw new ForbiddenException('Billing is disabled')
    }
    const organization = await this.organizationsService.findOne(orgname)

    // Cancel the subscription
    await this.billingService.cancelSubscription(organization.stripeCustomerId)
  }
}
